from __future__ import annotations

from abc import ABC, abstractmethod
from typing import Any, Optional


class LLMProvider(ABC):
    @abstractmethod
    def analyze_text(self, text: str, prompt: str, **kwargs: Any) -> str:
        ...

    @abstractmethod
    def extract_claims(self, text: str, **kwargs: Any) -> list[str]:
        ...

    @property
    @abstractmethod
    def name(self) -> str:
        ...

    @property
    @abstractmethod
    def available(self) -> bool:
        ...


class OpenAICompatibleProvider(LLMProvider):
    def __init__(self, api_key: str, base_url: str = "https://api.openai.com/v1", model: str = "gpt-4o"):
        self.api_key = api_key
        self.base_url = base_url.rstrip("/")
        self._model = model

    @property
    def name(self) -> str:
        return "OpenAI Compatible"

    @property
    def available(self) -> bool:
        return bool(self.api_key)

    def analyze_text(self, text: str, prompt: str, **kwargs: Any) -> str:
        if not self.available:
            return ""
        try:
            import httpx
            response = httpx.post(
                f"{self.base_url}/chat/completions",
                headers={
                    "Authorization": f"Bearer {self.api_key}",
                    "Content-Type": "application/json",
                },
                json={
                    "model": kwargs.get("model", self._model),
                    "messages": [
                        {"role": "system", "content": prompt},
                        {"role": "user", "content": text},
                    ],
                    "temperature": kwargs.get("temperature", 0.1),
                    "max_tokens": kwargs.get("max_tokens", 2048),
                },
                timeout=kwargs.get("timeout", 60),
            )
            response.raise_for_status()
            data = response.json()
            return data["choices"][0]["message"]["content"]
        except Exception:
            return ""

    def extract_claims(self, text: str, **kwargs: Any) -> list[str]:
        prompt = (
            "Extract all security-relevant declarative statements from the following text. "
            "Return each statement on a separate line. Only include statements that make "
            "a specific claim about security posture, access control, architecture, or process."
        )
        result = self.analyze_text(text, prompt, **kwargs)
        if not result:
            return []
        return [line.strip() for line in result.strip().split("\n") if line.strip()]


class OllamaProvider(LLMProvider):
    def __init__(self, base_url: str = "http://localhost:11434", model: str = "llama3"):
        self.base_url = base_url.rstrip("/")
        self._model = model

    @property
    def name(self) -> str:
        return "Ollama"

    @property
    def available(self) -> bool:
        try:
            import httpx
            response = httpx.get(f"{self.base_url}/api/tags", timeout=5)
            return response.is_success
        except Exception:
            return False

    def analyze_text(self, text: str, prompt: str, **kwargs: Any) -> str:
        if not self.available:
            return ""
        try:
            import httpx
            response = httpx.post(
                f"{self.base_url}/api/chat",
                json={
                    "model": kwargs.get("model", self._model),
                    "messages": [
                        {"role": "system", "content": prompt},
                        {"role": "user", "content": text},
                    ],
                    "options": {
                        "temperature": kwargs.get("temperature", 0.1),
                    },
                },
                timeout=kwargs.get("timeout", 120),
            )
            response.raise_for_status()
            data = response.json()
            return data.get("message", {}).get("content", "")
        except Exception:
            return ""

    def extract_claims(self, text: str, **kwargs: Any) -> list[str]:
        prompt = (
            "Extract all security-relevant declarative statements from the following text. "
            "Return each statement on a separate line. Only include statements that make "
            "a specific claim about security posture, access control, architecture, or process."
        )
        result = self.analyze_text(text, prompt, **kwargs)
        if not result:
            return []
        return [line.strip() for line in result.strip().split("\n") if line.strip()]
