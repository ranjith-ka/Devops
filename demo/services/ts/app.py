from flask import Flask, request, jsonify
import numpy as np
import pandas as pd
from datetime import datetime, timedelta
import json
from sklearn.ensemble import RandomForestRegressor, GradientBoostingRegressor
from sklearn.linear_model import LinearRegression
from sklearn.metrics import mean_absolute_error, mean_squared_error
from sklearn.preprocessing import StandardScaler
import warnings

warnings.filterwarnings("ignore")

app = Flask(__name__)

# Initialize models
models = {
    "linear": LinearRegression(),
    "random_forest": RandomForestRegressor(n_estimators=100, random_state=42),
    "gradient_boost": GradientBoostingRegressor(n_estimators=100, random_state=42),
}

scaler = StandardScaler()

@app.route("/health")
def health():
    return jsonify({"status": "ok", "models_available": list(models.keys())})


def create_features(df, target_col, lookback_window=7):
    """Create time series features"""
    features = df.copy()

    # Lag features
    for i in range(1, lookback_window + 1):
        features[f"lag_{i}"] = features[target_col].shift(i)

    # Rolling statistics
    for window in [3, 7, 14]:
        features[f"rolling_mean_{window}"] = (
            features[target_col].rolling(window=window).mean()
        )
        features[f"rolling_std_{window}"] = (
            features[target_col].rolling(window=window).std()
        )

    # Time-based features
    if "timestamp" in features.columns:
        features["timestamp"] = pd.to_datetime(features["timestamp"])
        features["day_of_week"] = features["timestamp"].dt.dayofweek
        features["month"] = features["timestamp"].dt.month
        features["quarter"] = features["timestamp"].dt.quarter
        features["is_weekend"] = (features["day_of_week"] >= 5).astype(int)

    # Drop rows with NaN values created by shifting
    features = features.dropna()

    return features


@app.route("/forecast", methods=["GET", "POST"])
def forecast():
    """Enhanced demand forecasting with multiple algorithms"""
    try:
        if request.method == "GET":
            # Return sample forecast for compatibility
            return jsonify({"forecast": [{"ts": "2025-09-13", "pred": 0.2}]})

        data = request.json

        # Parse input data
        historical_data = data.get("data", [])
        forecast_days = data.get("days", 30)
        model_type = data.get("model", "random_forest")
        confidence_interval = data.get("confidence_level", 0.95)

        if not historical_data:
            return jsonify({"error": "No historical data provided"}), 400

        if model_type not in models:
            return (
                jsonify(
                    {"error": f"Invalid model type. Available: {list(models.keys())}"}
                ),
                400,
            )

        # Convert to DataFrame
        df = pd.DataFrame(historical_data)

        # Ensure we have required columns
        if "value" not in df.columns:
            return jsonify({"error": "Data must contain 'value' column"}), 400

        # Add timestamp if not provided
        if "timestamp" not in df.columns:
            df["timestamp"] = pd.date_range(
                start=datetime.now() - timedelta(days=len(df) - 1),
                periods=len(df),
                freq="D",
            )

        # Create features
        features_df = create_features(df, "value")

        if len(features_df) < 10:
            return (
                jsonify(
                    {
                        "error": "Insufficient data for forecasting (minimum 10 points required)"
                    }
                ),
                400,
            )

        # Prepare training data
        feature_cols = [
            col for col in features_df.columns if col not in ["value", "timestamp"]
        ]
        X = features_df[feature_cols].fillna(0)
        y = features_df["value"]

        # Scale features
        X_scaled = scaler.fit_transform(X)

        # Split data for validation
        split_idx = int(len(X_scaled) * 0.8)
        X_train, X_val = X_scaled[:split_idx], X_scaled[split_idx:]
        y_train, y_val = y[:split_idx], y[split_idx:]

        # Train model
        model = models[model_type]
        model.fit(X_train, y_train)

        # Validate model
        val_predictions = model.predict(X_val)
        mae = mean_absolute_error(y_val, val_predictions)
        rmse = np.sqrt(mean_squared_error(y_val, val_predictions))

        # Generate forecasts
        forecasts = []
        last_values = features_df.tail(1).copy()

        for i in range(forecast_days):
            # Prepare features for next prediction
            feature_row = []

            # Use the same feature structure
            for col in feature_cols:
                if col.startswith("lag_"):
                    lag_num = int(col.split("_")[1])
                    if lag_num == 1:
                        feature_row.append(
                            last_values["value"].iloc[0]
                            if len(forecasts) == 0
                            else forecasts[-1]["value"]
                        )
                    else:
                        # Use historical data or previous forecasts
                        if lag_num <= len(forecasts):
                            feature_row.append(forecasts[-(lag_num - 1)]["value"])
                        else:
                            idx = len(features_df) - (lag_num - len(forecasts))
                            if idx >= 0:
                                feature_row.append(features_df["value"].iloc[idx])
                            else:
                                feature_row.append(0)
                else:
                    # For other features, use last known value or calculate
                    feature_row.append(
                        last_values[col].iloc[0] if col in last_values.columns else 0
                    )

            # Scale and predict
            feature_row_scaled = scaler.transform([feature_row])
            prediction = model.predict(feature_row_scaled)[0]

            # Add uncertainty bounds
            uncertainty = (
                mae * 1.96 if confidence_interval == 0.95 else mae * 2.58
            )  # 99% CI

            forecast_date = features_df["timestamp"].iloc[-1] + timedelta(days=i + 1)

            forecasts.append(
                {
                    "date": forecast_date.isoformat(),
                    "value": float(prediction),
                    "lower_bound": float(prediction - uncertainty),
                    "upper_bound": float(prediction + uncertainty),
                    "confidence": confidence_interval,
                }
            )

        # Calculate trend and seasonality insights
        recent_trend = np.polyfit(
            range(min(30, len(df))), df["value"].tail(min(30, len(df))), 1
        )[0]

        # Seasonal patterns (weekly)
        if len(df) >= 14:
            weekly_pattern = df.groupby(df.index % 7)["value"].mean().to_dict()
        else:
            weekly_pattern = {}

        return jsonify(
            {
                "forecasts": forecasts,
                "model_performance": {
                    "mae": float(mae),
                    "rmse": float(rmse),
                    "model_type": model_type,
                },
                "insights": {
                    "trend": (
                        "increasing"
                        if recent_trend > 0
                        else "decreasing" if recent_trend < 0 else "stable"
                    ),
                    "trend_slope": float(recent_trend),
                    "weekly_pattern": {
                        str(k): float(v) for k, v in weekly_pattern.items()
                    },
                    "forecast_period": f"{forecast_days} days",
                    "confidence_level": confidence_interval,
                },
                "processing_time": datetime.now().isoformat(),
            }
        )

    except Exception as e:
        return jsonify({"error": str(e)}), 500


@app.route("/analyze_patterns", methods=["POST"])
def analyze_patterns():
    """Analyze patterns in logistics data"""
    try:
        data = request.json
        historical_data = data.get("data", [])

        if not historical_data:
            return jsonify({"error": "No data provided"}), 400

        df = pd.DataFrame(historical_data)

        if "value" not in df.columns:
            return jsonify({"error": "Data must contain 'value' column"}), 400

        # Add timestamp if not provided
        if "timestamp" not in df.columns:
            df["timestamp"] = pd.date_range(
                start=datetime.now() - timedelta(days=len(df) - 1),
                periods=len(df),
                freq="D",
            )
        else:
            df["timestamp"] = pd.to_datetime(df["timestamp"])

        # Calculate statistics
        stats = {
            "mean": float(df["value"].mean()),
            "median": float(df["value"].median()),
            "std": float(df["value"].std()),
            "min": float(df["value"].min()),
            "max": float(df["value"].max()),
            "skewness": float(df["value"].skew()),
            "kurtosis": float(df["value"].kurtosis()),
        }

        # Detect anomalies using IQR method
        Q1 = df["value"].quantile(0.25)
        Q3 = df["value"].quantile(0.75)
        IQR = Q3 - Q1
        lower_bound = Q1 - 1.5 * IQR
        upper_bound = Q3 + 1.5 * IQR

        anomalies = df[(df["value"] < lower_bound) | (df["value"] > upper_bound)]

        # Seasonal analysis
        df["day_of_week"] = df["timestamp"].dt.dayofweek
        df["month"] = df["timestamp"].dt.month

        daily_patterns = (
            df.groupby("day_of_week")["value"].agg(["mean", "std"]).to_dict()
        )
        monthly_patterns = df.groupby("month")["value"].agg(["mean", "std"]).to_dict()

        # Trend analysis
        if len(df) >= 7:
            recent_trend = np.polyfit(range(len(df)), df["value"], 1)[0]
            trend_direction = (
                "increasing"
                if recent_trend > 0
                else "decreasing" if recent_trend < 0 else "stable"
            )
        else:
            recent_trend = 0
            trend_direction = "insufficient_data"

        # Volatility analysis
        rolling_std = df["value"].rolling(window=min(7, len(df))).std()
        volatility = "high" if rolling_std.mean() > stats["std"] else "low"

        return jsonify(
            {
                "statistics": stats,
                "anomalies": {
                    "count": len(anomalies),
                    "dates": anomalies["timestamp"].dt.strftime("%Y-%m-%d").tolist(),
                    "values": anomalies["value"].tolist(),
                    "threshold_lower": float(lower_bound),
                    "threshold_upper": float(upper_bound),
                },
                "patterns": {
                    "daily": {
                        str(k): {"mean": float(v)}
                        for k, v in daily_patterns["mean"].items()
                    },
                    "monthly": {
                        str(k): {"mean": float(v)}
                        for k, v in monthly_patterns["mean"].items()
                    },
                },
                "trend": {
                    "direction": trend_direction,
                    "slope": float(recent_trend),
                    "volatility": volatility,
                },
                "insights": {
                    "data_points": len(df),
                    "date_range": f"{df['timestamp'].min().strftime('%Y-%m-%d')} to {df['timestamp'].max().strftime('%Y-%m-%d')}",
                    "coefficient_of_variation": (
                        float(stats["std"] / stats["mean"]) if stats["mean"] != 0 else 0
                    ),
                },
            }
        )

    except Exception as e:
        return jsonify({"error": str(e)}), 500


@app.route("/capacity_planning", methods=["POST"])
def capacity_planning():
    """Capacity planning recommendations"""
    try:
        data = request.json

        current_capacity = data.get("current_capacity", 100)
        historical_demand = data.get("demand_data", [])
        growth_rate = data.get("expected_growth_rate", 0.1)  # 10% by default
        planning_horizon = data.get("planning_horizon_months", 12)

        if not historical_demand:
            return jsonify({"error": "No demand data provided"}), 400

        df = pd.DataFrame(historical_demand)
        if "value" not in df.columns:
            return jsonify({"error": "Data must contain 'value' column"}), 400

        # Calculate current utilization
        current_demand = df["value"].tail(30).mean()  # Average of last 30 days
        current_utilization = (current_demand / current_capacity) * 100

        # Project future demand
        monthly_projections = []
        base_demand = current_demand

        for month in range(1, planning_horizon + 1):
            projected_demand = base_demand * (1 + growth_rate) ** (month / 12)
            utilization = (projected_demand / current_capacity) * 100

            # Determine if capacity expansion is needed
            capacity_needed = projected_demand * 1.2  # 20% buffer
            expansion_needed = capacity_needed > current_capacity

            monthly_projections.append(
                {
                    "month": month,
                    "projected_demand": float(projected_demand),
                    "utilization_percent": float(utilization),
                    "capacity_needed": float(capacity_needed),
                    "expansion_needed": expansion_needed,
                    "additional_capacity": float(
                        max(0, capacity_needed - current_capacity)
                    ),
                }
            )

        # Find when capacity expansion is first needed
        expansion_timeline = next(
            (proj for proj in monthly_projections if proj["expansion_needed"]), None
        )

        # Recommendations
        recommendations = []
        if current_utilization > 90:
            recommendations.append(
                "Immediate capacity expansion recommended - current utilization > 90%"
            )
        elif current_utilization > 80:
            recommendations.append(
                "Monitor closely - utilization approaching capacity limits"
            )

        if expansion_timeline:
            recommendations.append(
                f"Plan capacity expansion by month {expansion_timeline['month']}"
            )

        # Cost analysis (simplified)
        peak_additional_capacity = max(
            [proj["additional_capacity"] for proj in monthly_projections]
        )
        estimated_expansion_cost = (
            peak_additional_capacity * 1000
        )  # $1000 per unit of capacity

        return jsonify(
            {
                "current_status": {
                    "capacity": current_capacity,
                    "current_demand": float(current_demand),
                    "utilization_percent": float(current_utilization),
                },
                "projections": monthly_projections,
                "recommendations": recommendations,
                "expansion_timeline": expansion_timeline,
                "cost_estimate": {
                    "peak_additional_capacity": float(peak_additional_capacity),
                    "estimated_cost": float(estimated_expansion_cost),
                    "currency": "USD",
                },
                "risk_assessment": {
                    "risk_level": (
                        "high"
                        if current_utilization > 85
                        else "medium" if current_utilization > 70 else "low"
                    ),
                    "buffer_available": float(current_capacity - current_demand),
                },
            }
        )

    except Exception as e:
        return jsonify({"error": str(e)}), 500


if __name__ == "__main__":
    app.run(host="0.0.0.0", port=8002, debug=True)
