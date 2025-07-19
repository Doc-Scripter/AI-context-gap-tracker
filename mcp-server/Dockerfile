# AI Context Gap Tracker - MCP Server Dockerfile
FROM python:3.11-slim

WORKDIR /app

# Install system dependencies
RUN apt-get update && apt-get install -y \
    build-essential \
    && rm -rf /var/lib/apt/lists/*

# Copy requirements and install Python dependencies
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Copy application code
COPY main.py .

# Set environment variables
ENV PYTHONPATH=/app
ENV TRACKER_API_ENDPOINT=http://ai-context-tracker:8080

# Expose port (for health checks or debugging)
EXPOSE 8000

# Set the entrypoint
ENTRYPOINT ["python", "main.py"]