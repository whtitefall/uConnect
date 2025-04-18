import faiss
import json
from sentence_transformers import SentenceTransformer
from llama_cpp import Llama

index = faiss.read_index("golang_code.index")
chunks = json.load(open("golang_chunks_meta.json"))
embed_model = SentenceTransformer("intfloat/e5-base")
llm = Llama(model_path="models/mistral-7b-instruct-v0.2.Q4_K_M.gguf", n_ctx=4096, n_gpu_layers=32)

def query_codebase(question: str, top_k: int = 4):
    q_embedding = embed_model.encode([question])
    D, I = index.search(q_embedding, top_k)
    selected = [chunks[i] for i in I[0]]
    return selected

def build_prompt(contexts, question):
    context_str = "\n\n".join([f"File: {c['file']}\nCode:\n{c['code']}" for c in contexts])
    return f"""You are an expert Golang engineer. Given the following code snippets and a user question, provide a helpful explanation.

Context:
{context_str}

Question: {question}
Answer:"""

while True:
    q = input("🧠 Ask about your Golang code: ")
    ctx = query_codebase(q)
    prompt = build_prompt(ctx, q)
    output = llm(prompt, max_tokens=512, stop=["User:", "Question:"], echo=False)
    print("🤖", output["choices"][0]["text"].strip())