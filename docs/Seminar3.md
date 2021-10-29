# Programarea aplicatiilor distribuite(PAD)
## Seminar 3
## Topic: API Architectural Styles

<hr>

### Students:  _Pavlov Alexandru_ , _Filipescu Mihai_, _Dodi Cristian-Dumitru_
### Group: __FAF-181__

<hr>

- What is the paper about?
- What is a server / client stub, in the context of RPC?
- What does it mean to be integrated with WS-security protocols? Exemplify some of these protocols and what they protect against.
- How do you understand HATEOAS?
- "GraphQL has subscriptions" - What are subscriptions? Why would we need them?
- Order the API patterns by message size.
- Which API pattern would best fit your laboratory work? Why?


### Answers

- This paper presents and compares 4 styles for an API architecture, including: REST, SOAP, GraphQL and RPC.
- Stub is the code which converts the arguments between the client and the server during procedures call.
- In the context of SOAP - WS secuirty protocols act as an extension for it in order to enforce secuirty. Main methods for ensuring security are signing SOAP messages, encription and simple tokens. These protocols do not provide complete security, and it is required to use other technologies as well. In the context of a web service, security is helping with eavesdropping, spoofing and backdoor attacks
- HATEOAS is the key part which distinguishes REST from other architectural styles. Main idea is that when a user interacts with a application, besides receiving normal responses, the server provides additional dynamic information. As a small example, after making a request to a REST, all possible future requests which may be called are given inside the response.
- Just like simple querries, subscriptions are a way to fetch data, with main difference being that the results can change over time. They are mostly used when fast updates are needed in real time.
- RPC, GraphQL, SOAP, REST
- We will use REST style, since we are the most familliar with this type of architecture, GraphQL will not be good for our nested structures for the game state, RPC is not ideal for an external API.
