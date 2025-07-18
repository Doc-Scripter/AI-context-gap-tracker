# AI-context-gap-tracker

## Project Summary
AIContextTracker is a modular system that augments language models with real-time context tracking, logic evaluation, and dynamic prompting to bridge information gaps in conversation. It is designed as a complement to existing LLMs and operates independently using well-defined logical rules, context validation layers, and follow-up intelligence. It does not merely ask clarifying questions; it asserts answers while noting ambiguities and then optionally seeks clarification.

## Motivation
Despite large language models’ impressive linguistic output, they often suffer from the following issues:
- Lack of real-time memory across conversation turns
- Assumption of user clarity even in ambiguous input
- Hallucinations due to partial data or unstated assumptions
- No module dedicated to truth-checking, logic rules, or assertion-based reasoning

AIContextTracker addresses this by:
1. Tracking conversational context across turns
2. Applying logical rules and factual boundaries
3. Identifying unknowns and ambiguous sections
4. Responding assertively, then marking and optionally requesting clarification

## How It Helps LLMs
LLMs, even with system prompts and memory, still struggle with:
- Gaps in multi-turn understanding
- Misinterpretation of incomplete context
- Lack of traceable logic

AIContextTracker provides:
- **Context Synchronization**: Maintains and updates a memory graph or tree with incoming context.
- **Response Accuracy Indexing**: Tags responses with clarity flags (e.g., assumed, inferred, ambiguous, verified).
- **Logic Validation Engine**: Applies domain-specific or global rules for consistent reasoning.
- **Prompt Tuning Relay**: Adjusts LLM input with clarity and context framing.

This improves:
- Output precision (estimated +20–35%)
- Reduction of hallucinated or misinformed replies
- User trust and feedback loops

## Architecture Overview
```
User Input --> Context Analyzer --> Rule Engine --> Response Preprocessor --> LLM
                                    |                            |
                           Ambiguity Detector            Clarity Tagger
```

## Modules
### 1. Context Tracker
- Tracks topics, terms, timelines, assertions
- Formats a memory graph (mutable)
- Uses entity resolution and dialog state tracking

### 2. Logic Rule Engine
- Custom rule DSL or pre-defined rules
- Evaluates conditions like:
  - Temporal consistency
  - Scope agreement
  - Missing information

### 3. Response Auditor
- Classifies responses by level of certainty
- Flags assumptions, contradictions
- Applies retry or flag mechanisms

### 4. Prompt Rewriter
- Wraps original user input + tracked context
- Adds disambiguation flags
- Optimizes output quality

## Stack Justification
### Primary Language: **Golang**
Chosen for:
- Fast backend performance
- Built-in concurrency (goroutines)
- Easy REST/gRPC APIs
- Clean code and fast debugging

### AI/NLP Companion: **Python**
Chosen for:
- Rich NLP ecosystem (spaCy, NLTK, Transformers)
- Ideal for logic pattern detection and vector embeddings
- Can serve AI services via REST or message queues

### Rust Considered, but Deferred
While Rust offers:
- Absolute memory control
- WASM deployment options
- High performance

It lacks:
- Integration ease with Python LLMs
- Fast prototyping capabilities

Rust may be revisited if performance bottlenecks arise in:
- Memory modeling of context graphs
- Embedding-heavy matrix operations
- Future edge deployment (WASM)

## Real-World Example Use Case
**Example:** A user says:
"I want to get a visa for a place I mentioned earlier, but I changed my mind."

LLMs might guess. AIContextTracker will:
1. Refer to past location references
2. Identify ambiguity: “which place?”
3. Inject clarity in LLM prompt:
```json
{"place": "unknown, multiple possible", "last_mentioned": "Canada", "action": "visa application"}
```
4. LLM response:
"Assuming you meant Canada. Please confirm if that’s still the case."

## Technologies
- Golang (core backend services)
- Python (NLP and AI integration layer)
- Redis (context memory cache)
- PostgreSQL (structured data rules)
- REST/gRPC (communication between modules)
- Docker (deployment)
- RabbitMQ/Kafka (event streaming)

## Future Plans
- Add WebSocket layer for real-time conversation pipelines
- WASM module rewrite in Rust (for edge inference)
- Graph-based visual memory explorer (UI)
- Rule DSL (domain-specific logic compiler)

## Naming Justification
**AIContextTracker** reflects:
- Direct focus on context tracking
- Real-time memory and state awareness
- Logic and clarity management, rather than ethical/moral reasoning

Previous name **AIConscience** was considered too abstract or morally framed.

## Status
Initial architectural planning complete. Implementation of Go backend modules and Python clarity checker in progress.

## Contributors
Lead: Clifford Otieno

## License
MIT License

