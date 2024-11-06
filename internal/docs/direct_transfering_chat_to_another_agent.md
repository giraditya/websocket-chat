```mermaid
    sequenceDiagram
        actor User
        actor MasterAgent
        actor Agent1
        actor Agent2
        Agent1 ->> MasterAgent: I want to transfer this chat Agent2
        MasterAgent -->> Agent1: Ok, i will transfer this chat
        MasterAgent ->> Agent2: Hi, I have a chat that was handed over to you from Agent1
        Agent1 -->> MasterAgent: Ok, i will take this chat
        Agent1 --> User: Response old chat
```