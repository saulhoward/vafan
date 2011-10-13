all:
		make -C pkg/vafan test install
		make -C cmd/vafan

clean:
		make -C pkg/vafan clean
		make -C cmd/vafan clean

.PHONY: all clean
