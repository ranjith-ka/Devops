from flask import Flask, request, jsonify
import re
import json
from datetime import datetime, timedelta
import spacy
from transformers import pipeline, AutoTokenizer, AutoModelForSequenceClassification

app = Flask(__name__)

# Load models
try:
    nlp = spacy.load("en_core_web_sm")
except:
    print(
        "Warning: spaCy model not found. Install with: python -m spacy download en_core_web_sm"
    )
    nlp = None

# Initialize sentiment analysis pipeline
sentiment_pipeline = pipeline(
    "sentiment-analysis", model="distilbert-base-uncased-finetuned-sst-2-english"
)

# Initialize classification pipeline for customs/logistics
classification_pipeline = pipeline(
    "text-classification", model="microsoft/DialoGPT-medium"
)

@app.route("/health")
def health():
    return jsonify({"status": "ok", "models_loaded": nlp is not None})

@app.route("/classify", methods=["POST"])
def classify():
    """Enhanced text classification with ML models"""
    data = request.json
    text = data.get("text", "")

    if not text:
        return jsonify({"error": "No text provided"}), 400

    # Basic keyword classification
    keywords = {
        "customs_issue": [
            "customs",
            "tariff",
            "duty",
            "hs code",
            "classification",
            "import",
            "export",
        ],
        "delay": ["delay", "late", "urgent", "asap", "expedite", "rush"],
        "hypercare": ["hypercare", "premium", "vip", "priority", "critical"],
        "documentation": ["document", "paperwork", "certificate", "permit", "license"],
        "shipping": ["shipping", "freight", "cargo", "vessel", "container"],
    }

    text_lower = text.lower()
    classifications = {}

    for category, words in keywords.items():
        score = sum(1 for word in words if word in text_lower)
        if score > 0:
            classifications[category] = score / len(words)

    # Get sentiment analysis
    sentiment_result = sentiment_pipeline(text[:512])  # Limit text length
    sentiment = sentiment_result[0]

    # Determine primary classification
    if classifications:
        primary_label = max(classifications.keys(), key=lambda k: classifications[k])
        confidence = min(classifications[primary_label] + 0.3, 1.0)
    else:
        primary_label = "general"
        confidence = 0.5

    # Extract entities if spaCy is available
    entities = []
    if nlp:
        doc = nlp(text)
        entities = [
            {
                "text": ent.text,
                "label": ent.label_,
                "start": ent.start_char,
                "end": ent.end_char,
            }
            for ent in doc.ents
        ]

    return jsonify(
        {
            "label": primary_label,
            "confidence": confidence,
            "all_classifications": classifications,
            "sentiment": {"label": sentiment["label"], "score": sentiment["score"]},
            "entities": entities,
        }
    )


@app.route("/extract_metadata", methods=["POST"])
def extract_metadata():
    """Extract structured metadata from logistics text"""
    data = request.json
    text = data.get("text", "")

    metadata = {
        "eta": None,
        "sla_days": None,
        "customer_tier": "standard",
        "urgency": "normal",
        "location": None,
        "container_numbers": [],
        "hs_codes": [],
    }

    # Extract dates for ETA
    date_patterns = [
        r"\d{4}-\d{2}-\d{2}",
        r"\d{2}/\d{2}/\d{4}",
        r"\d{1,2}\s+(jan|feb|mar|apr|may|jun|jul|aug|sep|oct|nov|dec)[a-z]*\s+\d{4}",
    ]

    for pattern in date_patterns:
        matches = re.findall(pattern, text, re.IGNORECASE)
        if matches:
            metadata["eta"] = matches[0]
            break

    # Extract container numbers
    container_pattern = r"[A-Z]{4}\d{7}"
    metadata["container_numbers"] = re.findall(container_pattern, text)

    # Extract HS codes
    hs_pattern = r"\b\d{4}[\.\s]?\d{2}[\.\s]?\d{2}\b"
    metadata["hs_codes"] = re.findall(hs_pattern, text)

    # Determine customer tier
    if any(word in text.lower() for word in ["vip", "premium", "gold", "platinum"]):
        metadata["customer_tier"] = "premium"
    elif any(word in text.lower() for word in ["enterprise", "corporate"]):
        metadata["customer_tier"] = "enterprise"

    # Determine urgency
    urgent_words = ["urgent", "asap", "rush", "emergency", "critical"]
    if any(word in text.lower() for word in urgent_words):
        metadata["urgency"] = "high"

    # Extract SLA from text
    sla_match = re.search(r"sla[:\s]+(\d+)\s*days?", text, re.IGNORECASE)
    if sla_match:
        metadata["sla_days"] = int(sla_match.group(1))

    return jsonify(metadata)


@app.route("/similarity", methods=["POST"])
def calculate_similarity():
    """Calculate text similarity for duplicate detection"""
    data = request.json
    text1 = data.get("text1", "")
    text2 = data.get("text2", "")

    if not text1 or not text2:
        return jsonify({"error": "Both text1 and text2 required"}), 400

    # Simple similarity calculation
    from difflib import SequenceMatcher

    # Normalize texts
    text1_norm = re.sub(r"[^\w\s]", "", text1.lower())
    text2_norm = re.sub(r"[^\w\s]", "", text2.lower())

    # Calculate similarity
    similarity = SequenceMatcher(None, text1_norm, text2_norm).ratio()

    # Word-based similarity
    words1 = set(text1_norm.split())
    words2 = set(text2_norm.split())
    word_similarity = (
        len(words1.intersection(words2)) / len(words1.union(words2))
        if words1.union(words2)
        else 0
    )

    # Combined score
    combined_score = (similarity + word_similarity) / 2

    return jsonify(
        {
            "similarity": similarity,
            "word_similarity": word_similarity,
            "combined_score": combined_score,
            "is_duplicate": combined_score > 0.8,
        }
    )


if __name__ == "__main__":
    app.run(host="0.0.0.0", port=8000, debug=True)
