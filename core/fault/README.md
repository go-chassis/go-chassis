#fault

The fault-injector is used for error (delay / abort) registration, and is generally used on the consumer side of microservices. scenes to be used:
>1. Unit test
>2. The simulation of server crashed;
>3. Analog network freeze

fault-injector supports two types of processing:
>1. delay: The middleware will delay for a specified time according to a certain probability;
>2. abort: The middleware reports an error according to the probability. If an error occurs, the entire link is interrupted.
 

Implementation principle:
>When the request arrives, for each strategy, an integer [0-100] randVal is randomly generated, and when percentage> = randVal + 1, the strategy is hit


Note:
> 1. When delay and abort are configured at the same time, delay is executed first, and then abort is executed;
> 2. The go-chassis framework supports the following protocols by default. If you need other protocols or customization, you need to register through InstallFaultInjectionPlugin.
        <br/>rest, <br/>rhighway, <br/>rdubbo
    
 
