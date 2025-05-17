# from pathlib import Path
# from docling.datamodel.base_models import InputFormat
# from docling.datamodel.pipeline_options import (
#     AcceleratorDevice,
#     AcceleratorOptions,
#     PdfPipelineOptions,
# )
# from docling.datamodel.settings import settings
# from docling.document_converter import DocumentConverter, PdfFormatOption
# from docling.datamodel.settings import settings
# from docling.document_converter import DocumentConverter, PdfFormatOption

# # def main():
# input_doc = Path("C:/Users/BHUVAN/Downloads/Documents/securityhub.pdf")

# # Explicitly set the accelerator
# # accelerator_options = AcceleratorOptions(
# #     num_threads=8, device=AcceleratorDevice.AUTO
# # )
# # accelerator_options = AcceleratorOptions(
# #     num_threads=8, device=AcceleratorDevice.CPU
# # )
# # accelerator_options = AcceleratorOptions(
# #     num_threads=8, device=AcceleratorDevice.MPS
# # )
# accelerator_options = AcceleratorOptions(
#     num_threads=32, device=AcceleratorDevice.CUDA
# )

# # easyocr doesnt support cuda:N allocation, defaults to cuda:0
# # accelerator_options = AcceleratorOptions(num_threads=8, device="cuda:1")

# pipeline_options = PdfPipelineOptions()
# pipeline_options.accelerator_options = accelerator_options
# pipeline_options.do_ocr = True
# pipeline_options.do_table_structure = True
# pipeline_options.table_structure_options.do_cell_matching = True

# converter = DocumentConverter(
#     format_options={
#         InputFormat.PDF: PdfFormatOption(
#             pipeline_options=pipeline_options,
#         )
#     }
# )

# # Enable the profiling to measure the time spent
# settings.debug.profile_pipeline_timings = True

# # Convert the document
# conversion_result = converter.convert(input_doc)
# doc = conversion_result.document

# # List with total time per document
# doc_conversion_secs = conversion_result.timings["pipeline_total"].times

# md = doc.export_to_markdown()
# # print(md)
# print(f"Conversion secs: {doc_conversion_secs}")

# # docs = loader.load()
# with open("doc.md", "w") as f:
#     # for doc in docs:
#     f.write(md)
from pathlib import Path
from docling.datamodel.base_models import InputFormat
from langchain.text_splitter import MarkdownTextSplitter # Import MarkdownTextSplitter
    # Required imports for create_markdown_faiss_index (expected to be in the global scope):
from langchain.docstore.document import Document
from langchain_huggingface import HuggingFaceEmbeddings
from langchain_community.vectorstores import Chroma
import os # Though not directly used in create_markdown_faiss_index, FAISS.save_local uses path operations.

# import document, faiss that's not imported

# from langchain

def create_markdown_faiss_index(
    markdown_content: str,
    index_path: str = "faiss_markdown_index",
    embeddings_model_name: str = "intfloat/e5-large-v2"
):
    """
    Chunks a Markdown string, creates embeddings, builds a FAISS index,
    and saves it locally.
    Relies on Document, HuggingFaceEmbeddings, and FAISS being available from imports.
    """
    print(f"\n--- Creating FAISS index from Markdown content at {index_path} ---")

    # 1. Create a Document object from the markdown content
    # Ensure Document is imported: from langchain.docstore.document import Document
    doc = Document(page_content=markdown_content, metadata={"source": "local_markdown"})

    # 2. Split the Markdown document
    markdown_splitter = MarkdownTextSplitter(chunk_size=1500, chunk_overlap=350)
    split_docs = markdown_splitter.split_documents([doc])
    print(f"Markdown content split into {len(split_docs)} chunks.")
    if not split_docs:
        print("No chunks were created from the markdown content. Aborting index creation.")
        return

    # 3. Initialize embeddings
    # Ensure HuggingFaceEmbeddings is imported: from langchain_huggingface import HuggingFaceEmbeddings
    try:
        embeddings = HuggingFaceEmbeddings(model_name=embeddings_model_name,model_kwargs={"device":"cuda"})
    except Exception as e:
        print(f"Error initializing HuggingFaceEmbeddings: {e}")
        print("Please ensure sentence-transformers is installed and the model is accessible.")
        return

    # 4. Create FAISS vector store
    # Ensure FAISS is imported: from langchain_community.vectorstores import FAISS
    try:
        print("Creating FAISS vector store...")
        # os.makedirs(index_path)
        vectorstore = Chroma.from_documents(split_docs, embeddings, persist_directory=index_path)
        print("FAISS vector store created successfully.")
    except Exception as e:
        print(f"Error creating FAISS vector store: {e}")
        print("Ensure 'faiss-cpu' or 'faiss-gpu' is installed.")
        return

    # 5. Save the FAISS index locally
    try:
        # Ensure os is imported if save_local relies on it for complex path ops,
        # though usually it handles paths directly.
        
        # vectorstore.save_local(index_path)
        print(f"FAISS index saved locally to folder: {index_path}")
    except Exception as e:
        print(f"Error saving FAISS index: {e}")

# Example usage of the new function:
# if __name__ == "__main__":
    # This __main__ block now focuses only on the FAISS index creation.
    # Ensure necessary global imports like Document, HuggingFaceEmbeddings, FAISS are present
    # in the script for this example to run.

    # Example: Create a FAISS index from a sample Markdown string
        
f = open("doc.md", "r")
sample_markdown = f.read()
f.close()
create_markdown_faiss_index(sample_markdown, index_path=r"E:\weird shit\Aventus\index")

    # You can then load this index later if needed:
    # print("\n--- Example: Attempting to load the created FAISS index ---")
    # try:
    #     # Ensure HuggingFaceEmbeddings and FAISS are imported globally
    #     embeddings_for_load = HuggingFaceEmbeddings(model_name="intfloat/e5-large-v2")
    #     # For newer versions of FAISS or Langchain, allow_dangerous_deserialization might be required
    #     # when loading pickle files. Use with caution if the index source is not trusted.
    #     loaded_vectorstore = FAISS.load_local(
    #         "my_markdown_faiss_index",
    #         embeddings_for_load,
    #         allow_dangerous_deserialization=True # Set to True if loading a pickled FAISS index
    #     )
    #     print("FAISS index 'my_markdown_faiss_index' loaded successfully.")
    #
    #     # Example query on the loaded index
    #     query = "What is the API_KEY for?"
    #     results = loaded_vectorstore.similarity_search(query, k=1)
    #     if results:
    #         print(f"Query: '{query}'")
    #         print(f"Result from loaded index: '{results[0].page_content}'")
    #     else:
    #         print(f"No results found for the query '{query}' in the loaded index.")
    # except ImportError as ie:
    #     print(f"Import error during loading: {ie}. Make sure all Langchain components are installed.")
    # except Exception as e:
    #     print(f"Error loading or querying FAISS index 'my_markdown_faiss_index': {e}")
    #     print("Ensure the index was created successfully in the 'my_markdown_faiss_index' folder.")
    #     print("Also, check that 'faiss-cpu' or 'faiss-gpu' is installed and HuggingFaceEmbeddings model is accessible.")