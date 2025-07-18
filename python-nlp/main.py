from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
from typing import List, Dict, Any, Optional
import spacy
import nltk
from nltk.sentiment import SentimentIntensityAnalyzer
from nltk.tokenize import sent_tokenize, word_tokenize
from nltk.corpus import stopwords
from nltk.stem import WordNetLemmatizer
import redis
import json
import os
import logging
from datetime import datetime
import re

# Configure logging
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

# Initialize FastAPI app
app = FastAPI(title="AI Context Gap Tracker NLP Service", version="1.0.0")

# Initialize NLP models
try:
    nlp = spacy.load("en_core_web_sm")
    logger.info("SpaCy model loaded successfully")
except OSError:
    logger.error("SpaCy model not found. Please install with: python -m spacy download en_core_web_sm")
    nlp = None

# Initialize NLTK
try:
    nltk.download('punkt', quiet=True)
    nltk.download('stopwords', quiet=True)
    nltk.download('wordnet', quiet=True)
    nltk.download('vader_lexicon', quiet=True)
    
    sia = SentimentIntensityAnalyzer()
    lemmatizer = WordNetLemmatizer()
    stop_words = set(stopwords.words('english'))
    logger.info("NLTK components initialized successfully")
except Exception as e:
    logger.error(f"NLTK initialization failed: {e}")
    sia = None
    lemmatizer = None
    stop_words = set()

# Initialize Redis
try:
    redis_host = os.getenv('REDIS_HOST', 'localhost')
    redis_port = int(os.getenv('REDIS_PORT', 6379))
    redis_client = redis.Redis(host=redis_host, port=redis_port, decode_responses=True)
    redis_client.ping()
    logger.info("Redis connected successfully")
except Exception as e:
    logger.error(f"Redis connection failed: {e}")
    redis_client = None

# Pydantic models
class TextInput(BaseModel):
    text: str
    session_id: Optional[str] = None
    turn_number: Optional[int] = None

class EntityResult(BaseModel):
    text: str
    label: str
    start: int
    end: int
    confidence: float

class TopicResult(BaseModel):
    topic: str
    confidence: float
    keywords: List[str]

class SentimentResult(BaseModel):
    compound: float
    positive: float
    negative: float
    neutral: float
    label: str

class AmbiguityResult(BaseModel):
    text: str
    type: str
    confidence: float
    suggestions: List[str]

class TimelineEvent(BaseModel):
    event: str
    timestamp: Optional[str] = None
    reference: str
    confidence: float

class NLPAnalysisResult(BaseModel):
    entities: List[EntityResult]
    topics: List[TopicResult]
    sentiment: SentimentResult
    ambiguities: List[AmbiguityResult]
    timeline_events: List[TimelineEvent]
    key_phrases: List[str]
    language: str
    readability_score: float
    processing_time: float

# Health check endpoint
@app.get("/health")
async def health_check():
    return {
        "status": "healthy",
        "service": "nlp-service",
        "models": {
            "spacy": nlp is not None,
            "nltk": sia is not None,
            "redis": redis_client is not None
        }
    }

# Entity extraction endpoint
@app.post("/entities", response_model=List[EntityResult])
async def extract_entities(input_data: TextInput):
    if not nlp:
        raise HTTPException(status_code=503, detail="SpaCy model not available")
    
    try:
        doc = nlp(input_data.text)
        entities = []
        
        for ent in doc.ents:
            entities.append(EntityResult(
                text=ent.text,
                label=ent.label_,
                start=ent.start_char,
                end=ent.end_char,
                confidence=0.8  # SpaCy doesn't provide confidence scores directly
            ))
        
        return entities
    except Exception as e:
        logger.error(f"Entity extraction failed: {e}")
        raise HTTPException(status_code=500, detail="Entity extraction failed")

# Topic extraction endpoint
@app.post("/topics", response_model=List[TopicResult])
async def extract_topics(input_data: TextInput):
    if not nlp:
        raise HTTPException(status_code=503, detail="SpaCy model not available")
    
    try:
        doc = nlp(input_data.text)
        topics = []
        
        # Extract noun phrases as potential topics
        noun_phrases = [chunk.text.lower() for chunk in doc.noun_chunks if len(chunk.text.strip()) > 2]
        
        # Extract named entities as topics
        entity_topics = [ent.text.lower() for ent in doc.ents if ent.label_ in ['PERSON', 'ORG', 'GPE', 'PRODUCT']]
        
        # Combine and deduplicate
        all_topics = list(set(noun_phrases + entity_topics))
        
        for topic in all_topics[:10]:  # Limit to top 10 topics
            topics.append(TopicResult(
                topic=topic,
                confidence=0.7,
                keywords=topic.split()
            ))
        
        return topics
    except Exception as e:
        logger.error(f"Topic extraction failed: {e}")
        raise HTTPException(status_code=500, detail="Topic extraction failed")

# Sentiment analysis endpoint
@app.post("/sentiment", response_model=SentimentResult)
async def analyze_sentiment(input_data: TextInput):
    if not sia:
        raise HTTPException(status_code=503, detail="NLTK sentiment analyzer not available")
    
    try:
        scores = sia.polarity_scores(input_data.text)
        
        # Determine sentiment label
        if scores['compound'] >= 0.05:
            label = "positive"
        elif scores['compound'] <= -0.05:
            label = "negative"
        else:
            label = "neutral"
        
        return SentimentResult(
            compound=scores['compound'],
            positive=scores['pos'],
            negative=scores['neg'],
            neutral=scores['neu'],
            label=label
        )
    except Exception as e:
        logger.error(f"Sentiment analysis failed: {e}")
        raise HTTPException(status_code=500, detail="Sentiment analysis failed")

# Ambiguity detection endpoint
@app.post("/ambiguities", response_model=List[AmbiguityResult])
async def detect_ambiguities(input_data: TextInput):
    try:
        ambiguities = []
        text = input_data.text.lower()
        
        # Detect ambiguous pronouns
        ambiguous_pronouns = ['it', 'this', 'that', 'they', 'them', 'he', 'she', 'him', 'her']
        for pronoun in ambiguous_pronouns:
            if re.search(r'\b' + pronoun + r'\b', text):
                ambiguities.append(AmbiguityResult(
                    text=pronoun,
                    type="ambiguous_pronoun",
                    confidence=0.7,
                    suggestions=[f"Specify what '{pronoun}' refers to"]
                ))
        
        # Detect vague quantifiers
        vague_quantifiers = ['some', 'many', 'few', 'several', 'most', 'a lot of']
        for quantifier in vague_quantifiers:
            if re.search(r'\b' + quantifier + r'\b', text):
                ambiguities.append(AmbiguityResult(
                    text=quantifier,
                    type="vague_quantifier",
                    confidence=0.6,
                    suggestions=[f"Specify a more precise quantity than '{quantifier}'"]
                ))
        
        # Detect temporal ambiguities
        temporal_vague = ['soon', 'later', 'recently', 'a while ago', 'sometime']
        for temporal in temporal_vague:
            if re.search(r'\b' + temporal + r'\b', text):
                ambiguities.append(AmbiguityResult(
                    text=temporal,
                    type="temporal_ambiguity",
                    confidence=0.8,
                    suggestions=[f"Specify a more precise time than '{temporal}'"]
                ))
        
        return ambiguities
    except Exception as e:
        logger.error(f"Ambiguity detection failed: {e}")
        raise HTTPException(status_code=500, detail="Ambiguity detection failed")

# Timeline event extraction endpoint
@app.post("/timeline", response_model=List[TimelineEvent])
async def extract_timeline_events(input_data: TextInput):
    if not nlp:
        raise HTTPException(status_code=503, detail="SpaCy model not available")
    
    try:
        doc = nlp(input_data.text)
        events = []
        
        # Extract time-related entities
        time_entities = [ent for ent in doc.ents if ent.label_ in ['DATE', 'TIME', 'EVENT']]
        
        for ent in time_entities:
            events.append(TimelineEvent(
                event=ent.text,
                timestamp=None,  # Would need more sophisticated parsing
                reference=ent.label_,
                confidence=0.7
            ))
        
        # Extract sentences with temporal keywords
        temporal_keywords = ['yesterday', 'today', 'tomorrow', 'next week', 'last month', 'ago', 'later', 'before', 'after']
        sentences = sent_tokenize(input_data.text)
        
        for sentence in sentences:
            sentence_lower = sentence.lower()
            for keyword in temporal_keywords:
                if keyword in sentence_lower:
                    events.append(TimelineEvent(
                        event=sentence.strip(),
                        timestamp=None,
                        reference=f"temporal_keyword_{keyword}",
                        confidence=0.6
                    ))
                    break
        
        return events
    except Exception as e:
        logger.error(f"Timeline extraction failed: {e}")
        raise HTTPException(status_code=500, detail="Timeline extraction failed")

# Key phrase extraction endpoint
@app.post("/keyphrases", response_model=List[str])
async def extract_key_phrases(input_data: TextInput):
    if not nlp:
        raise HTTPException(status_code=503, detail="SpaCy model not available")
    
    try:
        doc = nlp(input_data.text)
        key_phrases = []
        
        # Extract noun phrases
        for chunk in doc.noun_chunks:
            if len(chunk.text.strip()) > 2 and chunk.text.lower() not in stop_words:
                key_phrases.append(chunk.text.strip())
        
        # Extract named entities
        for ent in doc.ents:
            if ent.text not in key_phrases:
                key_phrases.append(ent.text)
        
        return key_phrases[:20]  # Return top 20 key phrases
    except Exception as e:
        logger.error(f"Key phrase extraction failed: {e}")
        raise HTTPException(status_code=500, detail="Key phrase extraction failed")

# Complete NLP analysis endpoint
@app.post("/analyze", response_model=NLPAnalysisResult)
async def analyze_text(input_data: TextInput):
    start_time = datetime.now()
    
    try:
        # Perform all NLP analyses
        entities = await extract_entities(input_data)
        topics = await extract_topics(input_data)
        sentiment = await analyze_sentiment(input_data)
        ambiguities = await detect_ambiguities(input_data)
        timeline_events = await extract_timeline_events(input_data)
        key_phrases = await extract_key_phrases(input_data)
        
        # Calculate processing time
        processing_time = (datetime.now() - start_time).total_seconds()
        
        # Calculate readability score (simplified)
        readability_score = calculate_readability_score(input_data.text)
        
        # Detect language
        language = detect_language(input_data.text)
        
        result = NLPAnalysisResult(
            entities=entities,
            topics=topics,
            sentiment=sentiment,
            ambiguities=ambiguities,
            timeline_events=timeline_events,
            key_phrases=key_phrases,
            language=language,
            readability_score=readability_score,
            processing_time=processing_time
        )
        
        # Cache result if Redis is available
        if redis_client and input_data.session_id:
            cache_key = f"nlp_analysis:{input_data.session_id}:{input_data.turn_number}"
            redis_client.setex(cache_key, 3600, json.dumps(result.dict()))
        
        return result
        
    except Exception as e:
        logger.error(f"Complete NLP analysis failed: {e}")
        raise HTTPException(status_code=500, detail="NLP analysis failed")

# Helper functions
def calculate_readability_score(text: str) -> float:
    """Calculate a simple readability score based on sentence and word length"""
    if not text.strip():
        return 0.0
    
    sentences = sent_tokenize(text)
    words = word_tokenize(text)
    
    if len(sentences) == 0 or len(words) == 0:
        return 0.0
    
    avg_sentence_length = len(words) / len(sentences)
    avg_word_length = sum(len(word) for word in words) / len(words)
    
    # Simple readability score (lower is more readable)
    readability = (avg_sentence_length * 0.5) + (avg_word_length * 0.3)
    
    # Normalize to 0-1 scale
    return min(1.0, max(0.0, 1.0 - (readability / 20)))

def detect_language(text: str) -> str:
    """Simple language detection - defaults to English"""
    # This is a placeholder - would use proper language detection library
    return "en"

# Cache management endpoints
@app.get("/cache/status")
async def cache_status():
    if not redis_client:
        return {"status": "unavailable"}
    
    try:
        info = redis_client.info()
        return {
            "status": "available",
            "connected_clients": info.get('connected_clients', 0),
            "used_memory": info.get('used_memory_human', '0B'),
            "total_keys": redis_client.dbsize()
        }
    except Exception as e:
        return {"status": "error", "message": str(e)}

@app.delete("/cache/clear")
async def clear_cache():
    if not redis_client:
        raise HTTPException(status_code=503, detail="Redis not available")
    
    try:
        redis_client.flushdb()
        return {"message": "Cache cleared successfully"}
    except Exception as e:
        logger.error(f"Cache clear failed: {e}")
        raise HTTPException(status_code=500, detail="Cache clear failed")

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=5000)