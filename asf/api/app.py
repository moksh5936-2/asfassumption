from __future__ import annotations

from pathlib import Path

from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from fastapi.responses import HTMLResponse, FileResponse
from fastapi.staticfiles import StaticFiles

from asf import __version__
from asf.analyzer import Analyzer
from asf.config import ASFConfig
from asf.settings import load_config

from .routes import router, init

HERE = Path(__file__).parent
STATIC_DIR = HERE / "static"


def create_app(config: ASFConfig | None = None) -> FastAPI:
    if config is None:
        config = load_config()

    app = FastAPI(
        title="ASF Validator API",
        version=__version__,
        description="Assumption Security Framework Validator — experimental research platform",
    )

    app.add_middleware(
        CORSMiddleware,
        allow_origins=["*"],
        allow_credentials=True,
        allow_methods=["*"],
        allow_headers=["*"],
    )

    db_instance = None
    analyzer_instance = None

    if config.db_path:
        from asf.db.database import Database
        db_instance = Database(config.db_path)

    analyzer_instance = Analyzer(config)
    init(db_instance, analyzer_instance, config)

    app.state.db = db_instance
    app.state.analyzer = analyzer_instance
    app.state.config = config

    app.include_router(router)

    if STATIC_DIR.exists():
        app.mount("/static", StaticFiles(directory=str(STATIC_DIR)), name="static")

    @app.get("/", response_class=HTMLResponse)
    async def root():
        html = (HERE / "static" / "index.html").read_text() if (HERE / "static" / "index.html").exists() else ""
        return HTMLResponse(html)

    @app.get("/health")
    async def health():
        return {"status": "healthy", "version": __version__}

    @app.on_event("shutdown")
    def shutdown():
        if db_instance:
            db_instance.close()
        if analyzer_instance:
            analyzer_instance.close()

    return app


app = create_app()
