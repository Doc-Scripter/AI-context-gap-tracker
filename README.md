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

## Getting Started

This project uses Docker and Docker Compose for easy setup and execution across different operating systems (Linux, WSL, Windows).

### Prerequisites

Before you begin, ensure you have the following installed on your system:

1.  **Git**: For cloning the project repository.
    *   **Linux/WSL**: `sudo apt update && sudo apt install git` (for Debian/Ubuntu-based systems)
    *   **Windows**: Download from [git-scm.com](https://git-scm.com/download/win)

2.  **Docker Desktop** (for Windows/macOS) or **Docker Engine & Docker Compose** (for Linux):
    *   **Windows**: Install [Docker Desktop for Windows](https://docs.docker.com/desktop/install/windows-install/). This includes Docker Engine and Docker Compose.
    *   **WSL**: Ensure you have [WSL 2 installed](https://docs.microsoft.com/en-us/windows/wsl/install) and enable WSL 2 integration in Docker Desktop settings (Settings -> Resources -> WSL Integration).
    *   **Linux**: Follow the official guides to install [Docker Engine](https://docs.docker.com/engine/install/) and [Docker Compose](https://docs.docker.com/compose/install/) for your distribution.

### Setup and Run

Follow these simple steps to get the AI Context Gap Tracker running:

1.  **Clone the Repository**:
    Open your terminal (or Git Bash/PowerShell on Windows) and clone the project to your local machine:
    ```bash
    git clone https://github.com/your-repo/AI-context-gap-tracker.git
    cd AI-context-gap-tracker
    ```
    *(Note: Replace the URL with the actual repository URL if it's different.)*

2.  **Configure Environment Variables**:
    The project uses an `.env.example` file for configuration. Copy this file to `.env` and modify it if you need to customize settings (e.g., API keys for external AI models, database connections).
    ```bash
    cp .env.example .env
    # Open the .env file with a text editor and make any necessary changes.
    ```

3.  **Build and Run with Docker Compose**:
    Navigate to the project's root directory in your terminal and run the following command. This will build the necessary Docker images and start all the services (Go backend, Python NLP service, Redis, etc.).
    ```bash
    docker-compose up --build
    ```
    *   The `--build` flag ensures that all service images are built from scratch, incorporating any recent changes.
    *   The services will run in the foreground, and you will see their logs in your terminal.

4.  **Verify Services (Optional)**:
    You can open a new terminal window and run `docker-compose ps` to confirm that all services are running (`Up` status).

5.  **Access the Application**:
    The main Go server will typically be accessible via `http://localhost:8080` (or another port specified in the `docker-compose.yml` or `.env` file). You can now interact with its API from your client applications or development tools.

That's it! The AI Context Gap Tracker is now running locally, ready to enhance your AI interactions.

## Contributors
Lead: Clifford Otieno

## License
MIT License

