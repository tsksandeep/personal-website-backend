import os
import sys
import tabula
import pandas as pd

from table_ocr import extract_tables, extract_cells, ocr_image, ocr_to_csv


def get_csv_from_excel(file_path: str, csv_file_path: str):
    data_xls = pd.read_excel(
        file_path, dtype=str, index_col=None)

    with open(csv_file_path, 'w+') as writer:
        data_xls.to_csv(writer, encoding='utf-8', index=False)


def get_csv_from_pdf(file_path: str, csv_file_path: str):
    data_pdf = tabula.read_pdf(
        file_path, guess=False, pages='all', stream=True)[0]

    with open(csv_file_path, 'w+') as writer:
        data_pdf.to_csv(writer, encoding='utf-8', index=False)


def get_csv_from_image(file_path: str, csv_file_path: str):
    _, tables = extract_tables.main([file_path])[0]
    ocr = [
        ocr_image.main(cell, None)
        for cell in extract_cells.main(tables[0])
    ]

    with open(csv_file_path, 'w+') as writer:
        writer.write(ocr_to_csv.text_files_to_csv(ocr))


def main(file_path: str):
    file_ext = os.path.splitext(file_path)[1]
    csv_file_path = os.path.splitext(file_path)[0] + ".csv"

    if (file_ext == ".xls" or file_ext == ".xlsx"):
        get_csv_from_excel(file_path=file_path,
                           csv_file_path=csv_file_path)
    elif (file_ext == ".png" or file_ext == ".jpg" or file_ext == ".jpeg"):
        get_csv_from_image(file_path=file_path,
                           csv_file_path=csv_file_path)
    elif file_ext == ".pdf":
        get_csv_from_pdf(file_path=file_path,
                         csv_file_path=csv_file_path)
    else:
        sys.tracebacklimit = None
        raise Exception("file extension not supported")

    print(csv_file_path)


if __name__ == "__main__":
    main(sys.argv[1])
