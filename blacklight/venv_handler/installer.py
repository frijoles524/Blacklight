import os
import re
import json
import importlib.util
from pip._internal import main as pip_main
from importlib.metadata import distribution, PackageNotFoundError

def get_global_env():
    return f"{os.path.abspath(__file__)}/../../.blacklight"

def is_module_installed(module_name):
    try:
        distribution(module_name)
        return True
    except PackageNotFoundError:
        return False

def install_dependencies(json_str, target_dir):
    try:
        data = json.loads(json_str)
    except json.JSONDecodeError as e:
        raise ValueError(f"Invalid JSON provided: {e}")
    installed = set()

    for dep in data:
        dep = dep.strip()
        if dep and dep not in installed:
            if is_module_installed(re.split(r"(==|>=|<=|>|<|~=)", dep, maxsplit=1)[0].strip()):
                print(f"Global installation found for {dep}. Skipping...")
                installed.add(dep)
                continue
            print(f"Installing {dep} into {target_dir}...")
            result = pip_main([
                "install",
                "--target", target_dir,
                dep
            ])
            if result != 0:
                raise RuntimeError(f"Failed to install {dep}, pip exit code {result}")
            installed.add(dep)

def install_dependency_global(dep):
    dep = dep.strip()
    if dep:
        if is_module_installed(re.split(r"(==|>=|<=|>|<|~=)", dep, maxsplit=1)[0].strip()):
            print(f"{dep} installed. Aborting installation.")
            return
        print(f"Installing global {dep}")
        result = pip_main([
            "install",
            "--target", get_global_env(),
            dep
        ])
        if result != 0:
            raise RuntimeError(f"Failed to install {dep}, pip exit code {result}")
