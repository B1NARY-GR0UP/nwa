namespace go hello

/**
 * Hello World service
 */
service HelloWorldService {
    /**
     * Returns a greeting message
     */
    string sayHello(1:string name)
}