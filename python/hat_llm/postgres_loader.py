from dataclasses import dataclass
from typing import Any


@dataclass(slots=True)
class CorpusRecord:
    kind: str
    title: str
    text: str
    metadata: dict[str, Any]


class PostgresCorpusLoader:
    def __init__(self, dsn: str) -> None:
        self.dsn = dsn

    def load_postmortems(self, specialist_id: str, limit: int = 100) -> list[CorpusRecord]:
        rows = self._query(
            """
            SELECT task_summary, what_went_wrong, root_cause, failure_classification, created_at
            FROM agent_core.memory_postmortems
            WHERE specialist_id = %s
            ORDER BY created_at DESC
            LIMIT %s
            """,
            (specialist_id, limit),
        )
        return [
            CorpusRecord(
                kind="postmortem",
                title=row[0],
                text=f"task={row[0]}\nwhat_went_wrong={row[1]}\nroot_cause={row[2]}",
                metadata={"failure_classification": row[3], "created_at": str(row[4])},
            )
            for row in rows
        ]

    def load_eval_cases(self, specialist_id: str, limit: int = 100) -> list[CorpusRecord]:
        rows = self._query(
            """
            SELECT category, case_title, prompt_input, expected_behavior, expected_output_or_criteria, created_at
            FROM agent_core.memory_eval_cases
            WHERE specialist_id = %s
            ORDER BY created_at DESC
            LIMIT %s
            """,
            (specialist_id, limit),
        )
        return [
            CorpusRecord(
                kind="eval_case",
                title=row[1],
                text=f"prompt={row[2]}\nexpected_behavior={row[3]}\nexpected_output_or_criteria={row[4]}",
                metadata={"category": row[0], "created_at": str(row[5])},
            )
            for row in rows
        ]

    def _query(self, sql: str, params: tuple[Any, ...]) -> list[tuple[Any, ...]]:
        try:
            import psycopg
        except ImportError as exc:  # pragma: no cover
            raise RuntimeError("psycopg is required for Postgres corpus loading; install hat-llm[postgres]") from exc

        with psycopg.connect(self.dsn) as conn:
            with conn.cursor() as cur:
                cur.execute(sql, params)
                return list(cur.fetchall())
