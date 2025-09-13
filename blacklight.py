import json
import os
import sys
import tempfile
import shutil
import zipfile
import urllib.request
import urllib.error
from pathlib import Path
from typing import Dict, List, Optional, Tuple, Any
from venv_handler import installer
from venv_handler import loader

class App:
    def __init__(self, name: str, version: str, url: str, entrypoint: str, dependencies: List[str]):
        self.name = name
        self.version = version
        self.url = url
        self.entrypoint = entrypoint
        self.dependencies = dependencies

    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> 'App':
        deps = data.get('dependencies', [])
        if isinstance(deps, str):
            try:
                deps = json.loads(deps)
            except json.JSONDecodeError:
                deps = [deps]
        elif not isinstance(deps, list):
            deps = [str(deps)]
        return cls(
            name=data['name'],
            version=data['version'],
            url=data['url'],
            entrypoint=data['entrypoint'],
            dependencies=[str(dep).strip() for dep in deps if dep]
        )

    def to_dict(self) -> Dict[str, Any]:
        return {
            'name': self.name,
            'version': self.version,
            'url': self.url,
            'entrypoint': self.entrypoint,
            'dependencies': self.dependencies
        }


class AppStore:
    def __init__(self):
        self.apps: Dict[str, Dict[str, App]] = {}

    def get_latest_version(self, name: str) -> str:
        if name not in self.apps or not self.apps[name]:
            raise ValueError(f"No versions found for app: {name}")

        versions = list(self.apps[name].keys())
        latest = versions[0]

        for version in versions[1:]:
            if self._compare_versions(version, latest) > 0:
                latest = version

        return latest

    def _compare_versions(self, v1: str, v2: str) -> int:
        try:
            parts1 = [int(x) for x in v1.split('.')]
            parts2 = [int(x) for x in v2.split('.')]

            max_len = max(len(parts1), len(parts2))
            parts1.extend([0] * (max_len - len(parts1)))
            parts2.extend([0] * (max_len - len(parts2)))

            for p1, p2 in zip(parts1, parts2):
                if p1 > p2:
                    return 1
                elif p1 < p2:
                    return -1
            return 0
        except ValueError:
            return (v1 > v2) - (v1 < v2)

    def get_entrypoint(self, name: str, version: str) -> Optional[str]:
        return self.apps.get(name, {}).get(version, None).entrypoint if version in self.apps.get(name, {}) else None

    def get_url(self, name: str, version: str) -> Optional[str]:
        return self.apps.get(name, {}).get(version, None).url if version in self.apps.get(name, {}) else None

    def get_app(self, name: str, version: str) -> Optional[App]:
        return self.apps.get(name, {}).get(version, None)


def load_apps(directory: str) -> AppStore:
    store = AppStore()
    json_files = Path(directory).glob("*.json")

    for json_file in json_files:
        try:
            with open(json_file, 'r', encoding='utf-8') as f:
                data = json.load(f)

            apps = [data] if isinstance(data, dict) else data
            for app_data in apps:
                app = App.from_dict(app_data)
                if app.name in store.apps:
                    original_name = app.name
                    app.name = f"{app.name}_{json_file.stem}"
                    if app.name in store.apps:
                        print(f"Duplicate app detected. {original_name} from {json_file.stem} has been skipped.", file=sys.stderr)
                        continue
                    print(f"Duplicate app detected. Renamed {original_name} {app.version} to {app.name}", file=sys.stderr)
                store.apps.setdefault(app.name, {})[app.version] = app

        except (json.JSONDecodeError, KeyError, FileNotFoundError) as e:
            print(f"Error loading {json_file}: {e}", file=sys.stderr)
            continue

    return store


def download_file(url: str, destination: str) -> None:
    try:
        with urllib.request.urlopen(url) as response:
            if response.getcode() != 200:
                raise urllib.error.HTTPError(url, response.getcode(),
                                             f"HTTP {response.getcode()}", None, None)
            with open(destination, 'wb') as f:
                shutil.copyfileobj(response, f)
    except urllib.error.URLError as e:
        raise RuntimeError(f"Failed to download {url}: {e}")


def extract_zip(zip_path: str, destination: str) -> None:
    with zipfile.ZipFile(zip_path, 'r') as zip_ref:
        for member in zip_ref.namelist():
            if os.path.isabs(member) or ".." in member:
                raise ValueError(f"Unsafe path in ZIP: {member}")

        temp_extract_dir = Path(destination) / "__temp_extract__"
        zip_ref.extractall(temp_extract_dir)

        extracted_items = list(temp_extract_dir.iterdir())
        if len(extracted_items) == 1 and extracted_items[0].is_dir():
            root_folder = extracted_items[0]
            for item in root_folder.iterdir():
                shutil.move(str(item), destination)
            shutil.rmtree(temp_extract_dir)
        else:
            for item in extracted_items:
                shutil.move(str(item), destination)
            shutil.rmtree(temp_extract_dir)


def install_software(app: App) -> None:
    target_dir = Path(f"{app.name}-{app.version}")

    if target_dir.exists():
        raise FileExistsError(f"App {app.name}-{app.version} already installed")

    with tempfile.NamedTemporaryFile(suffix='.zip', delete=False) as tmp_file:
        tmp_path = tmp_file.name

    try:
        print(f"Downloading {app.name} {app.version}...")
        download_file(app.url, tmp_path)

        print(f"Extracting to {target_dir}...")
        target_dir.mkdir(parents=True, exist_ok=True)
        extract_zip(tmp_path, str(target_dir))

        venv_path = target_dir / ".blacklight"
        print(f"Creating virtual environment folder...")
        venv_path.mkdir(parents=True, exist_ok=True)

        if app.dependencies:
            deps_json = json.dumps(app.dependencies)
            print(f"Installing dependencies: {deps_json}")
            installer.install_dependencies(deps_json, str(venv_path))

        print(f"Successfully installed {app.name} {app.version}")

    finally:
        if Path(tmp_path).exists():
            os.unlink(tmp_path)

def run_software(name: str, version: str, store: AppStore) -> None:
    target_dir = Path(f"{name}-{version}")

    if not target_dir.exists():
        raise FileNotFoundError(f"App {name}-{version} not installed")

    entrypoint = store.get_entrypoint(name, version)
    if not entrypoint:
        raise ValueError(f"No entrypoint found for {name} {version}")

    venv_path = target_dir / ".blacklight"
    entrypoint_path = target_dir / entrypoint

    if not entrypoint_path.exists():
        raise FileNotFoundError(f"Entrypoint not found: {entrypoint_path}")

    try:
        loader.load_site_packages(str(venv_path))

        print(f"Running {name} {version}...")

        with open(entrypoint_path, 'r', encoding='utf-8') as f:
            code = f.read()

        globals_dict = {
            '__name__': '__main__',
            '__file__': str(entrypoint_path),
        }

        original_path = sys.path.copy()
        sys.path.insert(0, str(target_dir))

        try:
            exec(code, globals_dict)
        finally:
            sys.path = original_path

    finally:
        loader.unload_site_packages(str(venv_path))