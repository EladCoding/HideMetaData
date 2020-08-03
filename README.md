# HideMetaData MixNet

A basic project that implement "MixNet", An anonymous messaging system between two end users.

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

### Prerequisites

```
go version 1.14.6
```

### Installing

Clone the project to your computer.

```
git clone https://github.com/EladCoding/HideMetaData.git
```

Compile the project.

```
cd $ProjectDir$
make
```

## Running the tests

Run an automatic test.
- simulate mixnet infraStructure, three servers, and two clients.
- send a few messages from the clients to the servers, through the mixnet.
- validate that all messages has sent successfully to the server, and that the server has sent a received message to the client.
- validate that the server has decrypted and read the message successfully.

```
$ProjectDir$\bin\$ExecutableFileName$ 1 [Round Slot Time]
```

Run statistics.
- calculate the throughput of the mixnet on this computer.
- calculate the goodput of the mixnet on this computer.
* Note that when simulating statistics, we run the whole net on a spesific network,
we estimate that real-time using will be about four times more effective than the results we get here.

```
$ProjectDir$\bin\$ExecutableFileName$ 2 [Round Slot Time]
```

Run Playing example.
- simulate mixnet infraStructure, three servers.
- simulate one client, for the user to play with.
- ask the user for a name of a server destination.
- ask the user for a message.
- send the message to the server, and ask for more messages.

```
$ProjectDir$\bin\$ExecutableFileName$ 3 [Round Slot Time]
```

## Deployment

On a server machine, run:

```
$ProjectDir$\bin\$ExecutableFileName$ 5 server $ServerName$
```

On a mediator machine, run:

```
$ProjectDir$\bin\$ExecutableFileName$ 5 mediator $MediatorName$ [Round Slot Time]
```


On a client machine, run:

```
$ProjectDir$\bin\$ExecutableFileName$ 5 client $ClientName$
```


*Note that before running this part you should update the specific mixnet architecture details, 
and implement the function that generate it. at:
$ProjectDir$\scripts\CreateNodesMap\CreateNodesMap function.

## Contributing

Yossi Gilad

## Authors

Elad Shoham
