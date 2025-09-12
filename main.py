from blacklight import load_apps, install_software, run_software

store = load_apps("lists")

app_name = "my_app"
latest_version = store.get_latest_version(app_name)
app = store.get_app(app_name, latest_version)

install_software(app)

run_software(app.name, app.version, store)
