import os
import sys
import signal
import platform
import argparse
import subprocess

import faiss
import json
from sentence_transformers import SentenceTransformer

index = faiss.read_index("golang_code.index")
chunks = json.load(open("golang_chunks_meta.json"))
embed_model = SentenceTransformer("intfloat/e5-base")

# def run_command(command, shell=False):
#     """Run a system command and ensure it succeeds."""
#     try:
#         subprocess.run(command, shell=shell, check=True)
#     except subprocess.CalledProcessError as e:
#         print(f"Error occurred while running command: {e}")
#         sys.exit(1)


def run_command(command, shell=False):
    try:
        result = subprocess.run(command, shell=shell, check=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE, text=True)
        print(result.stdout.strip())  # Only print the actual answer
    except subprocess.CalledProcessError as e:
        print(f"Error occurred while running command: {e}")
        sys.exit(1)


def run_inference(prompt):
    build_dir = "build"
    if platform.system() == "Windows":
        main_path = os.path.join(build_dir, "bin", "Release", "llama-cli.exe")
        if not os.path.exists(main_path):
            main_path = os.path.join(build_dir, "bin", "llama-cli")
    else:
        main_path = os.path.join(build_dir, "bin", "llama-cli")

    command = [
        f'{main_path}',
        '-m', args.model,
        '-n', str(args.n_predict),
        '-t', str(args.threads),
        '-p', prompt,
        '-ngl', '0',
        '-c', str(args.ctx_size),
        '--temp', str(args.temperature),
        "-b", "1",
        '--no-display-prompt'
    ]
    if args.conversation:
        command.append("-cnv")
    run_command(command)

def signal_handler(sig, frame):
    print("Ctrl+C pressed, exiting...")
    sys.exit(0)



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




if __name__ == "__main__":

    # Usage: python run_inference.py -p "Microsoft Corporation is an American multinational corporation and technology company headquartered in Redmond, Washington."


    parser = argparse.ArgumentParser(description='Run inference')
    parser.add_argument("-m", "--model", type=str, help="Path to model file", required=False, default="models/BitNet-b1.58-2B-4T/ggml-model-i2_s.gguf")
    parser.add_argument("-n", "--n-predict", type=int, help="Number of tokens to predict when generating text", required=False, default=128)
    parser.add_argument("-p", "--prompt", type=str, help="Prompt to generate text from", required=False, default=None)
    parser.add_argument("-t", "--threads", type=int, help="Number of threads to use", required=False, default=2)
    parser.add_argument("-c", "--ctx-size", type=int, help="Size of the prompt context", required=False, default=2048)
    parser.add_argument("-temp", "--temperature", type=float, help="Temperature, a hyperparameter that controls the randomness of the generated text", required=False, default=0.8)
    parser.add_argument("-cnv", "--conversation", action='store_true', help="Whether to enable chat mode or not (for instruct models.)")


    args = parser.parse_args()
    
    while True:
        signal.signal(signal.SIGINT, signal_handler)
        q = input("🧠 Ask about your golang code: ")
        ctx = query_codebase(q)
        prompt = build_prompt(ctx, q)

        # args.prompt = prompt  # update prompt dynamically
        run_inference(prompt)

