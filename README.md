# HideMetaData MixNet

A basic project that implement "MixNet", An anonymous messaging system between two end users.

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

### Prerequisites

```
go version 1.14.6
```

### Installing

Clone the project to your computer

```
git clone https://github.com/EladCoding/HideMetaData.git
```

Compile the project

```
cd $ProjectDir$
make
```

## Running the tests

Run an automatic test. the test:
- simulate mixnet infraStructure, three servers, and two clients.
- send a few messages from the clients to the servers, through the mixnet.
- validate that all messages has sent successfully to the server, and that the server has sent a received message to the client.
- validate that the server has decrypted and read the message successfully.

```
$ProjectDir$\bin\HideMetaData 1
```

Run statistics.
- calculate the throughput of the mixnet on this computer.
- calculate the goodput of the mixnet on this computer.
* Note that when simulating statistics, we run the whole net on a spesific network,
we estimate that real-time using will be about four times more effective than the results we get here.

```
$ProjectDir$\bin\HideMetaData 2
```

Run Playing example.
- simulate mixnet infraStructure, three servers.
- simulate one client, for the user to play with.
- ask the user for a name of a server destination.
- ask the user for a message.
- send the message to the server, and ask for more messages.

```
$ProjectDir$\bin\HideMetaData 3
```

## Deployment

On a server machine, run:

```
$ProjectDir$\bin\HideMetaData 5 server $ServerName$
```

On a mediator machine, run:

```
$ProjectDir$\bin\HideMetaData 5 mediator $MediatorName$
```
*Note that before running this part you should update the specific mixnet architecture details, at:
$ProjectDir$\externals

## Contributing

Yossi Gilad

## Authors

Elad Shoham
