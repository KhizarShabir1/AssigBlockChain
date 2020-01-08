# Assigment2IBC
#.........................(BlockChain)..................
We have X+1 different nodes. The first being satoshi, the creator of genesis block and the seed for all other nodes. All other nodes connect to Satoshi first. We can organize the tasks to be done from the perspective of Satoshi and the Others as Shown below:

1. Satoshi:
On startup, the Satoshi node starts listening on a port specified using command line. It mints the first block, pays himself 100 coins (name your currency) and waits for other to connect.

It does not start serving unless a fixed amount of nodes, X, are connected, again specified using the command line. Satoshi mints one more block whenever a new node connects and hence becomes a rich node when all X have connected.

2. Others:
On startup, all other nodes start listening on an address (hostname:portNumber) provided as command line argument. This address will be used by other peers to connect to them. More on this later.

They all then connect with Satoshi (address provided using the command line) and tell Satoshi about the address they are listening at. Satoshi saves this information and tells them to wait, till the quorum is complete, that is, X nodes have are connected to Satoshi.

3. Satoshi:
When all X nodes have connected with Satoshi, it sends them the chain he has with X blocks all belonging to Satoshi. It also sends each of them information about the peers they need to connect to (along with their addresses) with the logic below:

The first node (in connection order) that connected to Satoshi is given one random node to connect to. The second node that connected to Satoshi is given two random nodes to connect to. The Xth node that connected to Satoshi is given all X random nodes to connect to.

So Lets say nodes are Alice, Bob, Joe, Charlie, John and Jim ,that is X=6, and they all are connected to Satoshi in the order given. Alice gets one random node to connect to, Bob gets two, Joe three and Jim (from the TrollHunter) is given the addresses of all other nodes to.

4. Others:
Once all other nodes get the chain and list of nodes and addresses to connect to, from Satoshi, they connect to them. At this point, all the nodes have current chain, all of them are connected to some random nodes and the setup is complete. Nothing happens then, and the network remains idle.

5. Satoshi and all other nodes:
Each node is then stuck at user input (in a separate goroutine and thus not interfering with other work a node is doing), asking nodes if they want to do some transactions.

For the start, Satoshi has some coins and no other node has any. So satoshi makes a transaction to send some coins to any user say Alice. Satoshi then randomly selects another node say Bob, to be the Miner for this transaction.

The transaction is only propagated to Bob then. Bob checks the validity of transaction, i.e. if Satoshi has indeed these coins. It then traverses the chain and calculate Satoshi’s balance. If the transaction is valid, Bob creates a new block with two transaction one provided and one for himself (as mining reward) and then floods the new block directly or indirectly to all other nodes. All other nodes then validate the block and update their chain.

The chain and peers sleep, till the next transaction !!

To further illustrate the process, lets say now Bob creates a new transaction (user input) and sends some coins to Jim. The transaction is sent to Satoshi, which picks a random node, say Charlie. Forwards transaction to Charlie, who validates, mines and propagates the new block.

The chain and peers sleep, till the next transaction !!
