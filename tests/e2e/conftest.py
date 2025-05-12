import logging
import pytest
from collections.abc import MutableMapping
from typing import Any

logging.basicConfig(
    filename='tavern_requests.log',
    level=logging.INFO,
    format='%(asctime)s - %(levelname)s - %(message)s',
    encoding='utf-8',
    filemode='w'
)

# Пример фикстуры, которая "заплаткает" отсутствующий атрибут
@pytest.fixture(autouse=True)
def add_global_cfg(request):
    if hasattr(request.node, 'global_cfg'):
        return
    # Например, считываем глобальные настройки из pytest.ini
    config = request.config.getini("tavern-global-cfg")
    # Создаём объект, содержащий переменные
    class GlobalCfg:
        def __init__(self, cfg):
            self.variables = {}  # или заполните по файлу cfg, если необходимо

    request.node.global_cfg = GlobalCfg(config)

def pytest_addoption(parser):
    # Регистрируем новые ini-опции, чтобы pytest не выдавал предупреждения
    parser.addini(
        "tavern-global-cfg",
        "Global config for Tavern tests",
        default=[],    # по умолчанию пустой список
        type="linelist"  # разбиваем значение на строки, чтобы вернуть список
    )
    parser.addini("tavern-strict", "Strict mode for Tavern tests", default="false")

@pytest.hookimpl(optionalhook=True)
def pytest_tavern_beta_before_every_request(request_args: MutableMapping):
    message = f"Request: {request_args['method']} {request_args['url']}"
    params = request_args.get('params')
    if params:
        message += f"\nQuery parameters: {params}"

    message += f"\nRequest body: {request_args.get('json', '<no body>')}"
    logging.info(message)

@pytest.hookimpl(optionalhook=True)
def pytest_tavern_beta_after_every_response(expected: Any, response: Any) -> None:
    logging.info(f"Response: {response.status_code} {response.text}")