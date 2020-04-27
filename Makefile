# Go parameters
GOCMD=go
GORELEASE=releases
GOPACKAGENAME=integrity-checker
GOCHECKFILE=$(GORELEASE)/checkfile.sha256
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOINSTALL=$(GOCMD) install
BINARY_NAME=_x64_$(GOPACKAGENAME)

all: test build-all install

build-all: clean
	GOOS=linux $(GOBUILD) -o $(GORELEASE)/linux$(BINARY_NAME).bin -v
	GOOS=darwin $(GOBUILD) -o $(GORELEASE)/darwin$(BINARY_NAME).bin -v
	GOOS=windows $(GOBUILD) -o $(GORELEASE)/windows$(BINARY_NAME).exe -v
	shasum -a 256 $(GORELEASE)/linux$(BINARY_NAME).bin $(GORELEASE)/darwin$(BINARY_NAME).bin $(GORELEASE)/windows$(BINARY_NAME).exe > $(GOCHECKFILE)

test:
	$(GOTEST) -v ./...

clean:
	$(GOCLEAN)
	rm -f $(GOCHECKFILE)
	rm -f $(GORELEASE)/linux$(BINARY_NAME).bin
	rm -f $(GORELEASE)/darwin$(BINARY_NAME).bin
	rm -f $(GORELEASE)/windows$(BINARY_NAME).exe

install:
	$(GOCLEAN) -i
	$(GOINSTALL) .