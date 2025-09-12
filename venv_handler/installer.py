import json
from pip._internal import main as pip_main

def install_dependencies(json_str, target_dir):
    try:
        data = json.loads(json_str)
    except json.JSONDecodeError as e:
        raise ValueError(f"Invalid JSON provided: {e}")
    installed = set()

    for dep in data:
        dep = dep.strip()
        if dep and dep not in installed:
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
        print(f"Installing global {dep}")
        result = pip_main([
            "install",
            dep
        ])
        if result != 0:
            raise RuntimeError(f"Failed to install {dep}, pip exit code {result}")
