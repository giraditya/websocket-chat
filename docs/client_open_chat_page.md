```mermaid
    sequenceDiagram
        actor User
        actor MasterAgent
        actor Customer Support Agent
        MasterAgent ->> User: Give the FAQ
        alt Is available answer from FAQ
            User -->> MasterAgent: Thanks!
        else I need customer support
            MasterAgent -->> Customer Support Agent: Notify the customer support for take this chat
            Customer Support Agent -->> MasterAgent: I take this chat
        end
        Customer Support Agent -->> User: Hi im Customer Support, can i help you? 
```