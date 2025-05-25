from docxtpl import DocxTemplate
import io
from grpc_clients.table_client import get_table_data
import logging

logger = logging.getLogger(__name__)


def generate_docx_from_template(template_bytes: bytes,
                                institution_id: int) -> bytes:
    template_stream = io.BytesIO(template_bytes)
    doc = DocxTemplate(template_stream)

    table_data = get_table_data(institution_id)

    logger.info(f"Generating docx for institution {institution_id}")
    logger.debug(f"Table columns: {table_data['columns']}")
    if table_data["rows"]:
        logger.debug(f"First row: {table_data['rows'][0]}")
    else:
        logger.warning("No rows received from gRPC")

    # table_data = {
    #     "columns": ["Column 1", "Column 2"],
    #     "rows": [["value 1", "value 2"], ["value 3", "value 4"]]
    #     }
    subdoc = doc.new_subdoc()
    table = subdoc.add_table(rows=1, cols=len(table_data["columns"]))
    try:
        table.style = 'Table Grid'
    except Exception:
        logger.exception("No table style")
        pass

    hdr_cells = table.rows[0].cells
    for i, col_name in enumerate(table_data["columns"]):
        hdr_cells[i].text = col_name

    for row in table_data["rows"]:
        row_cells = table.add_row().cells
        for i in range(min(len(row), len(row_cells))):
            row_cells[i].text = row[i]

    for i in range(len(row), len(row_cells)):
        row_cells[i].text = ""

    context = {
        "table": subdoc
    }

    doc.render(context)
    output_stream = io.BytesIO()
    doc.save(output_stream)
    return output_stream.getvalue()
