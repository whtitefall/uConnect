import json
import os
from sentence_transformers import SentenceTransformer
import faiss
import numpy as np

chunks = []
with open("golang_chunks.jsonl", "r", encoding="utf-8", errors="replace") as f:
    for i, line in enumerate(f):
        try:
            chunk = json.loads(line)
            if "code" in chunk:
                chunks.append(chunk)
        except json.JSONDecodeError as e:
            print(f"[警告] 第 {i+1} 行无法解析为 JSON，已跳过：{e}")


            
model = SentenceTransformer("intfloat/e5-base")  # or "BAAI/bge-base-en"
texts = [chunk["code"] for chunk in chunks]
embeddings = model.encode(texts, show_progress_bar=True)

dim = embeddings[0].shape[0]
index = faiss.IndexFlatL2(dim)
index.add(np.array(embeddings))

faiss.write_index(index, "golang_code.index")



with open("golang_chunks_meta.json", "w") as f:
    json.dump(chunks, f)