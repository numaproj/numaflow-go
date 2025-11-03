# OnSuccess Sink Example

An example User Defined Sink that simulates primary sink write failures/successes and returns a 
* Fallback response on write failures so that the payload will be sent to the fallback sink.
* OnSuccess response on write successes so that the payload will be sent to the onSuccess sink.