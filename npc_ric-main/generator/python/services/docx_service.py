from docx import Document
import io
from grpc_clients.table_client import get_table_data


def generate_docx_from_template(
        template_bytes: bytes,
        institution_id: int) -> bytes:
    template_stream = io.BytesIO(template_bytes)
    doc = Document(template_stream)

    table_data = get_table_data(institution_id)

    table = doc.add_table(rows=1, cols=len(table_data["columns"]))
    table.style = 'Table Grid'

    hdr_cells = table.rows[0].cells
    for i, col_name in enumerate(table_data["columns"]):
        hdr_cells[i].text = col_name

    for row in table_data["rows"]:
        row_cells = table.add_row().cells
        for i, cell_value in enumerate(row):
            row_cells[i].text = cell_value

    output_stream = io.BytesIO()
    doc.save(output_stream)
    return output_stream.getvalue()
