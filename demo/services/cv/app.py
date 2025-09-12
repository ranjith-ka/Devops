from flask import Flask, request, jsonify
import cv2
import numpy as np
import base64
import io
from PIL import Image
import easyocr
import torch
from ultralytics import YOLO
import json
from datetime import datetime

app = Flask(__name__)

# Initialize models
try:
    # OCR Reader
    ocr_reader = easyocr.Reader(["en"], gpu=torch.cuda.is_available())

    # Load YOLO model for object detection
    yolo_model = YOLO("yolov8n.pt")  # Download will happen on first run

    models_loaded = True
except Exception as e:
    print(f"Error loading models: {e}")
    ocr_reader = None
    yolo_model = None
    models_loaded = False

@app.route("/health")
def health():
    return jsonify(
        {
            "status": "ok",
            "models_loaded": models_loaded,
            "cuda_available": torch.cuda.is_available(),
        }
    )


def decode_image(image_data):
    """Decode base64 image data"""
    try:
        if "," in image_data:
            image_data = image_data.split(",")[1]

        image_bytes = base64.b64decode(image_data)
        image = Image.open(io.BytesIO(image_bytes))
        return cv2.cvtColor(np.array(image), cv2.COLOR_RGB2BGR)
    except Exception as e:
        raise ValueError(f"Invalid image data: {e}")


@app.route("/extract_text", methods=["POST"])
def extract_text():
    """Enhanced OCR text extraction with document type detection"""
    if not ocr_reader:
        return jsonify({"error": "OCR model not loaded"}), 500

    try:
        data = request.json
        image_data = data.get("image")

        if not image_data:
            return jsonify({"error": "No image provided"}), 400

        # Decode image
        img = decode_image(image_data)

        # Perform OCR
        results = ocr_reader.readtext(img)

        # Extract text and confidence scores
        extracted_text = []
        full_text = ""

        for bbox, text, confidence in results:
            if confidence > 0.5:  # Filter low confidence results
                extracted_text.append(
                    {"text": text, "confidence": float(confidence), "bbox": bbox}
                )
                full_text += text + " "

        # Document type detection based on keywords
        doc_type = "unknown"
        full_text_lower = full_text.lower()

        if any(
            word in full_text_lower for word in ["bill of lading", "bl", "shipping"]
        ):
            doc_type = "bill_of_lading"
        elif any(word in full_text_lower for word in ["invoice", "total", "amount"]):
            doc_type = "invoice"
        elif any(
            word in full_text_lower for word in ["customs", "declaration", "tariff"]
        ):
            doc_type = "customs_declaration"
        elif any(
            word in full_text_lower for word in ["certificate", "origin", "quality"]
        ):
            doc_type = "certificate"
        elif any(word in full_text_lower for word in ["packing", "list", "items"]):
            doc_type = "packing_list"

        # Extract structured data based on document type
        structured_data = extract_structured_data(full_text, doc_type)

        return jsonify(
            {
                "text": full_text.strip(),
                "document_type": doc_type,
                "structured_data": structured_data,
                "detailed_results": extracted_text,
                "processing_time": datetime.now().isoformat(),
            }
        )

    except Exception as e:
        return jsonify({"error": str(e)}), 500


def extract_structured_data(text, doc_type):
    """Extract structured data based on document type"""
    import re

    data = {}

    if doc_type == "bill_of_lading":
        # Extract container numbers
        container_pattern = r"[A-Z]{4}\d{7}"
        data["container_numbers"] = re.findall(container_pattern, text)

        # Extract vessel name
        vessel_match = re.search(r"vessel[:\s]+([^\n]+)", text, re.IGNORECASE)
        if vessel_match:
            data["vessel"] = vessel_match.group(1).strip()

        # Extract port of loading/discharge
        pol_match = re.search(r"port\s+of\s+loading[:\s]+([^\n]+)", text, re.IGNORECASE)
        if pol_match:
            data["port_of_loading"] = pol_match.group(1).strip()

    elif doc_type == "invoice":
        # Extract total amount
        amount_pattern = r"total[:\s]+\$?([0-9,]+\.?\d*)"
        amount_match = re.search(amount_pattern, text, re.IGNORECASE)
        if amount_match:
            data["total_amount"] = amount_match.group(1)

        # Extract invoice number
        invoice_pattern = r"invoice\s+(?:number|no)[:\s#]*([A-Z0-9\-]+)"
        invoice_match = re.search(invoice_pattern, text, re.IGNORECASE)
        if invoice_match:
            data["invoice_number"] = invoice_match.group(1)

    elif doc_type == "customs_declaration":
        # Extract HS codes
        hs_pattern = r"\b\d{4}[\.\s]?\d{2}[\.\s]?\d{2}\b"
        data["hs_codes"] = re.findall(hs_pattern, text)

        # Extract declared value
        value_pattern = r"declared\s+value[:\s]+\$?([0-9,]+\.?\d*)"
        value_match = re.search(value_pattern, text, re.IGNORECASE)
        if value_match:
            data["declared_value"] = value_match.group(1)

    return data


@app.route("/detect_objects", methods=["POST"])
def detect_objects():
    """Object detection for containers, vehicles, and cargo"""
    if not yolo_model:
        return jsonify({"error": "YOLO model not loaded"}), 500

    try:
        data = request.json
        image_data = data.get("image")

        if not image_data:
            return jsonify({"error": "No image provided"}), 400

        # Decode image
        img = decode_image(image_data)

        # Run YOLO detection
        results = yolo_model(img)

        detections = []
        for r in results:
            boxes = r.boxes
            if boxes is not None:
                for box in boxes:
                    class_id = int(box.cls[0])
                    confidence = float(box.conf[0])
                    bbox = box.xyxy[0].tolist()

                    detections.append(
                        {
                            "class_name": yolo_model.names[class_id],
                            "confidence": confidence,
                            "bbox": bbox,
                            "class_id": class_id,
                        }
                    )

        # Filter for logistics-relevant objects
        relevant_objects = [
            "truck",
            "car",
            "bus",
            "train",
            "boat",
            "airplane",
            "person",
            "motorcycle",
            "bicycle",
            "traffic light",
        ]

        logistics_detections = [
            det
            for det in detections
            if det["class_name"] in relevant_objects and det["confidence"] > 0.5
        ]

        # Count objects by type
        object_counts = {}
        for det in logistics_detections:
            obj_type = det["class_name"]
            object_counts[obj_type] = object_counts.get(obj_type, 0) + 1

        return jsonify(
            {
                "total_objects": len(logistics_detections),
                "object_counts": object_counts,
                "detections": logistics_detections,
                "processing_time": datetime.now().isoformat(),
            }
        )

    except Exception as e:
        return jsonify({"error": str(e)}), 500


@app.route("/detect", methods=["POST"])
def detect():
    """Legacy endpoint - redirects to detect_objects"""
    return detect_objects()


@app.route("/quality_check", methods=["POST"])
def quality_check():
    """Quality assessment of documents and cargo images"""
    try:
        data = request.json
        image_data = data.get("image")
        check_type = data.get("type", "document")  # document or cargo

        if not image_data:
            return jsonify({"error": "No image provided"}), 400

        # Decode image
        img = decode_image(image_data)

        quality_score = 0.0
        issues = []

        # Calculate image quality metrics
        gray = cv2.cvtColor(img, cv2.COLOR_BGR2GRAY)

        # 1. Brightness check
        mean_brightness = np.mean(gray)
        if mean_brightness < 50:
            issues.append("Image too dark")
            quality_score -= 0.2
        elif mean_brightness > 200:
            issues.append("Image too bright")
            quality_score -= 0.1
        else:
            quality_score += 0.3

        # 2. Blur detection using Laplacian variance
        laplacian_var = cv2.Laplacian(gray, cv2.CV_64F).var()
        if laplacian_var < 100:
            issues.append("Image is blurry")
            quality_score -= 0.3
        else:
            quality_score += 0.3

        # 3. Contrast check
        contrast = gray.std()
        if contrast < 30:
            issues.append("Low contrast")
            quality_score -= 0.2
        else:
            quality_score += 0.2

        # 4. Resolution check
        height, width = gray.shape
        if width < 800 or height < 600:
            issues.append("Low resolution")
            quality_score -= 0.2
        else:
            quality_score += 0.2

        # Normalize quality score
        quality_score = max(0.0, min(1.0, quality_score + 0.5))

        # Determine overall quality
        if quality_score > 0.8:
            overall_quality = "excellent"
        elif quality_score > 0.6:
            overall_quality = "good"
        elif quality_score > 0.4:
            overall_quality = "fair"
        else:
            overall_quality = "poor"

        recommendations = []
        if "too dark" in str(issues):
            recommendations.append("Increase lighting or camera exposure")
        if "blurry" in str(issues):
            recommendations.append("Hold camera steady and ensure proper focus")
        if "Low contrast" in str(issues):
            recommendations.append("Improve lighting conditions")
        if "Low resolution" in str(issues):
            recommendations.append("Use higher resolution camera or move closer")

        return jsonify(
            {
                "quality_score": round(quality_score, 2),
                "overall_quality": overall_quality,
                "issues": issues,
                "recommendations": recommendations,
                "metrics": {
                    "brightness": round(mean_brightness, 2),
                    "blur_score": round(laplacian_var, 2),
                    "contrast": round(contrast, 2),
                    "resolution": f"{width}x{height}",
                },
            }
        )

    except Exception as e:
        return jsonify({"error": str(e)}), 500


if __name__ == "__main__":
    app.run(host="0.0.0.0", port=8001, debug=True)
