# all the imports
from __future__ import with_statement
from contextlib import closing
from flask import Flask, request, session, g, redirect, url_for, \
     abort, render_template, Markup, flash
import datetime
import os
import yaml
import markdown
import smartypants
import settings
import hashlib
import sqlite3

# create the app...
app = Flask(__name__)
app.config.update(
    SECRET_KEY=settings.SECRET_KEY
)

def connect_db():
    return sqlite3.connect(app.config['DATABASE'])

def init_db():
    with closing(connect_db()) as db:
        with app.open_resource('schema.sql') as f:
            db.cursor().executescript(f.read())
        db.commit()
        with app.open_resource('initial-data.sql') as f:
            db.cursor().executescript(f.read())
        db.commit()

@app.before_request
def before_request():
    """Make sure we are connected to the database each request."""
    g.db = connect_db()

@app.after_request
def after_request(response):
    """Closes the database again at the end of the request."""
    g.db.close()
    return response

def query_db(query, args=(), one=False):
    cur = g.db.execute(query, args)
    rv = [dict((cur.description[idx][0], value)
               for idx, value in enumerate(row)) for row in cur.fetchall()]
    return (rv[0] if rv else None) if one else rv

def format_datetime(datetime_object, format):
    """Format a datetime object for display, used in Jinja2 templates"""
    return datetime_object.strftime(format)

def get_page(directory, file):
    """Load and parse a page from the filesystem. Returns the page, or None if not found"""
    path = os.path.abspath(os.path.join(os.path.dirname(__file__), directory, str(file) + '.mkd'))
    try:
        file_contents = open(path).read()
        file_contents = unicode(file_contents, 'utf-8')
    except:
        return None

    page_data = file_contents.split('>>>\n', 2)
    data = page_data[0] 
    text = page_data[1]
    sections = None
    if len(page_data) > 2:
        sections = page_data[2]
        sections = sections.split('>>>\n')
        sections = map(convert_markdown_to_html, sections)

    page = yaml.load(data)
    page['content'] = convert_markdown_to_html(text)
    page['path'] = file
    if sections is not None:
        page['sections'] = sections
    return page

def convert_markdown_to_html(text):
    text = markdown.markdown(text)
    text = smartypants.smartyPants(text)
    return Markup(text)

# Views 
@app.route('/')
def home():
    # page = get_page(settings.PAGES_DIR, 'home')
    # if page is None:
        # abort(404)
    return render_template('home.html')

@app.route('/<path>/')
def page(path):
    page = get_page(settings.PAGES_DIR, path)
    if page is None:
        abort(404)
    return render_template('page.html', page=page)

@app.errorhandler(404)
def page_not_found(error):
    return render_template('404.html'), 404

# Add jinja filters
app.jinja_env.filters['datetimeformat'] = format_datetime

if __name__ == '__main__':
    app.run(host='0.0.0.0',debug=settings.DEBUG)
