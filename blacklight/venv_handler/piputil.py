import urllib.request
import importlib.util
import types

def pip_installed():
    return importlib.util.find_spec('pip') is not None

def load_get_pip():
    with urllib.request.urlopen("https://bootstrap.pypa.io/get-pip.py") as response:
        source_code = response.read().decode('utf-8')
    module = types.ModuleType("_getpip_module")
    exec(source_code, module.__dict__)

    def wrapper(*args, **kwargs):
        if "__main__" in module.__dict__:
            module.__dict__["__name__"] = "__main__"
        if "main" in module.__dict__:
            return module.main(*args, **kwargs)
        else:
            return None

    wrapper.__name__ = "run_pip_installer"
    return wrapper

def install_pip():
    if not pip_installed():
        print("Pip was not found. Installing...")
        installer = load_get_pip()
        installer()