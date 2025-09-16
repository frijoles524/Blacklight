import urllib.request
import importlib.util
import types
import ast

def pip_installed():
    return importlib.util.find_spec('pip') is not None

def load_get_pip():
    try:
        with urllib.request.urlopen("https://bootstrap.pypa.io/get-pip.py") as response:
            source_code = response.read().decode("utf-8")
    except Exception as e:
        print("Error while loading pip installer:", e)
        raise Exception("Unable to load pip installer") from e

    source_code = '__name__ = "__main__"\n' + source_code # trick get_pip into running
    tree = ast.parse(source_code, filename="get-pip.py")

    wrapper_body = [ast.FunctionDef(
        name="_getpip_main",
        args=ast.arguments(
            posonlyargs=[], args=[], kwonlyargs=[], kw_defaults=[], defaults=[]
        ),
        body=tree.body,
        decorator_list=[]
    )]
    ast.fix_missing_locations(ast.Module(body=wrapper_body, type_ignores=[]))

    code = compile(ast.Module(body=wrapper_body, type_ignores=[]), filename="get-pip.py", mode="exec")

    module = types.ModuleType("_getpip_module")
    exec(code, module.__dict__)

    def run_pip_installer():
        return module._getpip_main()

    run_pip_installer.__name__ = "run_pip_installer"
    return run_pip_installer

def install_pip():
    if not pip_installed():
        print("Pip was not found. Installing...")
        load_get_pip()()