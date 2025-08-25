# Simple Source

A simple example of a user-defined source.
The SimpleSource is a basic source implementation that listens for changes in a side
input file using a file watcher. When the specified side input file is created,
its content is read, stored in a global variable, and sent through a global channel.
This content is then consumed by the source and sent to the associated message channel.

The source maintains an array of messages and implements the `Read()`, `Ack()`, and `Pending()` methods.

The `Read(x)` function of the SimpleSource struct is responsible for reading data from the source (in this case, from the globalChan channel, which receives file content when a side input file is created) and sending this data to a specified message channel.
The `Ack()` method acknowledges the last batch of messages returned by `Read()`.
The `Pending()` method returns 0 to indicate that the simple source always has 0 pending messages.