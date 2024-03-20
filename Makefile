BINDIR=${CURDIR}/bin
PACKAGE=${CURDIR}/cmd/app/

bindir:
	mkdir -p ${BINDIR}

build: bindir
	go build -o ${BINDIR}/app ${PACKAGE}