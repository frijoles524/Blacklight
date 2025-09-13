import sys
import site
from pathlib import Path

def load_site_packages(path: str):
    abs_path = str(Path(path).resolve())
    if abs_path not in sys.path:
        sys.path.insert(0, abs_path)
        site.addsitedir(abs_path)

def unload_site_packages(path: str):
    abs_path = str(Path(path).resolve())
    if abs_path in sys.path:
        sys.path = [p for p in sys.path if p != abs_path]
    to_delete = [
        name for name, module in sys.modules.items()
        if getattr(module, "__file__", None) and abs_path in str(module.__file__)
    ]
    for name in to_delete:
        del sys.modules[name]