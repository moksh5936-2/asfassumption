"""Generate a PDF version of the sample finance policy for testing."""
from pathlib import Path

try:
    from fpdf import FPDF
except ImportError:
    import subprocess, sys
    subprocess.check_call([sys.executable, "-m", "pip", "install", "fpdf2"])
    from fpdf import FPDF

SAMPLE_DIR = Path(__file__).parent.parent / "sample_data"
TXT_PATH = SAMPLE_DIR / "finance_policy.txt"
PDF_PATH = SAMPLE_DIR / "finance_policy.pdf"


def generate():
    text = TXT_PATH.read_text()
    lines = text.split("\n")

    pdf = FPDF()
    pdf.add_page()
    pdf.set_font("Helvetica", "B", 16)
    pdf.cell(0, 10, "Finance Access Control Policy", new_x="LMARGIN", new_y="NEXT")
    pdf.ln(4)

    for line in lines:
        line = line.strip()
        if not line:
            pdf.ln(3)
            continue
        if line.isupper() and len(line) > 3:
            pdf.set_font("Helvetica", "B", 11)
            pdf.cell(0, 8, line, new_x="LMARGIN", new_y="NEXT")
            pdf.set_font("Helvetica", "", 10)
        else:
            safe = line.encode("latin-1", "replace").decode("latin-1")
            pdf.multi_cell(0, 5, safe, new_x="LMARGIN", new_y="NEXT")

    pdf.output(str(PDF_PATH))
    print(f"Generated {PDF_PATH} ({PDF_PATH.stat().st_size} bytes)")


if __name__ == "__main__":
    generate()
