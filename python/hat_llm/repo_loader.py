from pathlib import Path

from .types import SlotBundle


class RepositoryLoader:
    def __init__(self, repo_root: str | Path = ".") -> None:
        self.repo_root = Path(repo_root)

    def load_specialist_slots(self, specialist_id: str) -> SlotBundle:
        base = self.repo_root / "specialists" / specialist_id
        return SlotBundle(
            identity=self._read(base / "IDENTITY.md"),
            soul=self._read(base / "SOUL.md"),
            agents=self._read(base / "AGENTS.md"),
            tools=self._read(base / "TOOLS.md"),
            skills=self._read(base / "SKILLS.md"),
            prompt=self._read(base / "PROMPT.md"),
            heartbeat=self._read(base / "HEARTBEAT.md"),
            memory=self._read(base / "MEMORY.md"),
            dreams=self._read(base / "DREAMS.md"),
            postmortem=self._read(base / "POSTMORTEM.md"),
        )

    def _read(self, path: Path) -> str:
        if not path.exists():
            return ""
        return path.read_text(encoding="utf-8")
