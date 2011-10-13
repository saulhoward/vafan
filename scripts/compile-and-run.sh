rm pkg/vafan/_go_.8
rm cmd/vafan/_go_.8
gomake install -C pkg/vafan/
gomake install -C cmd/vafan/
/usr/lib/go/bin/vafan
