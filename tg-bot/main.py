import asyncio
import logging
import requests
from aiogram import Bot, Dispatcher
from aiogram.filters import Command
from aiogram.types import Message
import os
from dotenv import load_dotenv
from aiogram.client.default import DefaultBotProperties
from aiogram.types import ReplyKeyboardMarkup, KeyboardButton, InlineKeyboardMarkup, InlineKeyboardButton
load_dotenv()

# Вставьте сюда токен вашего бота
TELEGRAM_BOT_TOKEN = os.getenv('TG_TOKEN')

# Базовый URL вашего API
API_BASE_URL = os.getenv('API_URL')

# Включение логирования
logging.basicConfig(level=logging.INFO)

# Инициализация бота и диспетчера
bot = Bot(token=TELEGRAM_BOT_TOKEN , default=DefaultBotProperties(parse_mode="markdown"))
dp = Dispatcher()

# Обработчик команды /start

# Keyboards
main_menu = ReplyKeyboardMarkup(
    keyboard=[
        [
            KeyboardButton(text='/clients'),
            KeyboardButton(text='/advertisers'),
        ],
        [
            KeyboardButton(text='/campaigns'),
            KeyboardButton(text='/ads'),
        ],
        [
            KeyboardButton(text='/stats'),
            KeyboardButton(text='/time'),
        ],
    ],
    resize_keyboard=True
)
@dp.message(Command(commands=['start']))
async def cmd_start(message: Message):
    help_text = (
        "Добро пожаловать! Я могу помочь вам с следующими командами:\n\n"
        "- `/get_client <client_id>`: Получить информацию о клиенте по его ID.\n"
        "- `/get_advertiser <advertiser_id>`: Получить информацию о рекламодателе по его ID.\n"
        "- `/get_campaigns <advertiser_id> [page] [size]`: Получить список кампаний рекламодателя с пагинацией.\n"
        "- `/get_campaign <advertiser_id> <campaign_id>`: Получить информацию о кампании по её ID.\n"
        "- `/get_ad <client_id>`: Получить рекламное объявление для клиента.\n"
        "- `/get_campaign_stats <campaign_id>`: Получить статистику по кампании.\n"
        "- `/get_advertiser_stats <advertiser_id>`: Получить агрегированную статистику по всем кампаниям рекламодателя.\n"
        "- `/get_campaign_daily_stats <campaign_id>`: Получить ежедневную статистику по кампании.\n"
        "- `/get_advertiser_daily_stats <advertiser_id>`: Получить ежедневную агрегированную статистику по кампаниям рекламодателя.\n\n"
        "Используйте команды, чтобы получить нужную информацию."
    )
    await message.reply(help_text,reply_markup=main_menu)

# Clients Menu
@dp.message(Command(commands=['clients']))
async def clients_menu(message: Message):
    keyboard = InlineKeyboardMarkup(inline_keyboard=[
        [
            InlineKeyboardButton(text='Get Client by ID', callback_data='get_client'),
        ],
        [
            InlineKeyboardButton(text='Upsert Clients (Bulk)', callback_data='upsert_clients'),
        ],
    ])
    await message.reply("Клиенты:", reply_markup=keyboard)

# Advertisers Menu
@dp.message(Command(commands=['advertisers']))
async def advertisers_menu(message: Message):
    keyboard = InlineKeyboardMarkup(inline_keyboard=[
        [
            InlineKeyboardButton(text='Get Advertiser by ID', callback_data='get_advertiser'),
        ],
        [
            InlineKeyboardButton(text='Upsert Advertisers (Bulk)', callback_data='upsert_advertisers'),
        ],
    ])
    await message.reply("Рекламодатели:", reply_markup=keyboard)

# Campaigns Menu
@dp.message(Command(commands=['campaigns']))
async def campaigns_menu(message: Message):
    keyboard = InlineKeyboardMarkup(inline_keyboard=[
        [
            InlineKeyboardButton(text='Create Campaign', callback_data='create_campaign'),
        ],
        [
            InlineKeyboardButton(text='Get Campaigns', callback_data='get_campaigns'),
        ],
        [
            InlineKeyboardButton(text='Get Campaign by ID', callback_data='get_campaign_by_id'),
        ],
        [
            InlineKeyboardButton(text='Update Campaign', callback_data='update_campaign'),
        ],
        [
            InlineKeyboardButton(text='Delete Campaign', callback_data='delete_campaign'),
        ],
    ])
    await message.reply("Кампании:", reply_markup=keyboard)

# Ads Menu
@dp.message(Command(commands=['ads']))
async def ads_menu(message: Message):
    keyboard = InlineKeyboardMarkup(inline_keyboard=[
        [
            InlineKeyboardButton(text='Get Ad for Client', callback_data='get_ad_for_client'),
        ],
        [
            InlineKeyboardButton(text='Record Ad Click', callback_data='record_ad_click'),
        ],
    ])
    await message.reply("Рекламные объявления:", reply_markup=keyboard)

# Stats Menu
@dp.message(Command(commands=['stats']))
async def stats_menu(message: Message):
    keyboard = InlineKeyboardMarkup(inline_keyboard=[
        [
            InlineKeyboardButton(text='Campaign Stats', callback_data='campaign_stats'),
        ],
        [
            InlineKeyboardButton(text='Advertiser Campaign Stats', callback_data='advertiser_campaign_stats'),
        ],
        [
            InlineKeyboardButton(text='Campaign Daily Stats', callback_data='campaign_daily_stats'),
        ],
        [
            InlineKeyboardButton(text='Advertiser Daily Stats', callback_data='advertiser_daily_stats'),
        ],
    ])
    await message.reply("Статистика:", reply_markup=keyboard)

# Time Menu
@dp.message(Command(commands=['time']))
async def time_menu(message: Message):
    keyboard = InlineKeyboardMarkup(inline_keyboard=[
        [
            InlineKeyboardButton(text='Advance Time', callback_data='advance_time'),
        ],
    ])
    await message.reply("Управление временем:", reply_markup=keyboard)

# Callback Query Handlers
@dp.callback_query()
async def handle_callbacks(query):
    if query.data == 'get_client':
        await query.message.reply("Введите ID клиента (UUID):")
        @dp.message()
        async def process_get_client(message: Message):
            client_id = message.text.strip()
            url = f"{API_BASE_URL}/clients/{client_id}"
            response = requests.get(url)
            if response.status_code == 200:
                client = response.json()
                reply_message = (
                    f"**Клиент ID**: `{client['client_id']}`\n"
                    f"**Логин**: `{client['login']}`\n"
                    f"**Возраст**: `{client['age']}`\n"
                    f"**Локация**: `{client['location']}`\n"
                    f"**Пол**: `{client['gender']}`"
                )
                await message.reply(reply_message)
            else:
                await message.reply(f"Ошибка: {response.status_code}")
            dp.message_handlers.unregister(process_get_client)

    elif query.data == 'upsert_clients':
        await query.message.reply("Отправьте данные клиентов в формате JSON:")
        @dp.message()
        async def process_upsert_clients(message: Message):
            try:
                data = message.text.strip()
                clients = requests.post(f"{API_BASE_URL}/clients/bulk", json=data)
                if clients.status_code == 201:
                    await message.reply("Клиенты успешно добавлены/обновлены.")
                else:
                    await message.reply(f"Ошибка: {clients.status_code}")
            except Exception as e:
                await message.reply(f"Ошибка: {e}")
            dp.message_handlers.unregister(process_upsert_clients)

    elif query.data == 'get_advertiser':
        await query.message.reply("Введите ID рекламодателя (UUID):")
        @dp.message()
        async def process_get_advertiser(message: Message):
            advertiser_id = message.text.strip()
            url = f"{API_BASE_URL}/advertisers/{advertiser_id}"
            response = requests.get(url)
            if response.status_code == 200:
                advertiser = response.json()
                reply_message = (
                    f"**Рекламодатель ID**: `{advertiser['advertiser_id']}`\n"
                    f"**Название**: `{advertiser['name']}`"
                )
                await message.reply(reply_message)
            else:
                await message.reply(f"Ошибка: {response.status_code}")
            dp.message_handlers.unregister(process_get_advertiser)

    elif query.data == 'upsert_advertisers':
        await query.message.reply("Отправьте данные рекламодателей в формате JSON:")
        @dp.message()
        async def process_upsert_advertisers(message: Message):
            try:
                data = message.text.strip()
                advertisers = requests.post(f"{API_BASE_URL}/advertisers/bulk", json=data)
                if advertisers.status_code == 201:
                    await message.reply("Рекламодатели успешно добавлены/обновлены.")
                else:
                    await message.reply(f"Ошибка: {advertisers.status_code}")
            except Exception as e:
                await message.reply(f"Ошибка: {e}")
            dp.message_handlers.unregister(process_upsert_advertisers)

    elif query.data == 'create_campaign':
        await query.message.reply("Отправьте данные кампании в формате JSON вместе с advertiser_id в URL:")
        @dp.message()
        async def process_create_campaign(message: Message):
            try:
                data = message.text.strip()
                advertiser_id = data.get('advertiser_id')
                url = f"{API_BASE_URL}/advertisers/{advertiser_id}/campaigns"
                campaign = requests.post(url, json=data)
                if campaign.status_code == 201:
                    await message.reply("Кампания успешно создана.")
                else:
                    await message.reply(f"Ошибка: {campaign.status_code}")
            except Exception as e:
                await message.reply(f"Ошибка: {e}")
            dp.message_handlers.unregister(process_create_campaign)

    elif query.data == 'get_campaigns':
        await query.message.reply("Введите advertiser_id:")
        @dp.message()
        async def process_get_campaigns(message: Message):
            advertiser_id = message.text.strip()
            url = f"{API_BASE_URL}/advertisers/{advertiser_id}/campaigns"
            response = requests.get(url)
            if response.status_code == 200:
                campaigns = response.json()
                if not campaigns:
                    await message.reply("Кампании не найдены.")
                    return
                reply_message = "Список кампаний:\n"
                for campaign in campaigns:
                    reply_message += (
                        f"-**Кампания ID**: `{campaign['campaign_id']}`, "
                        f"**Название**: `{campaign['ad_title']}`\n"
                    )
                await message.reply(reply_message)
            else:
                await message.reply(f"Ошибка: {response.status_code}")
            dp.message_handlers.unregister(process_get_campaigns)

    elif query.data == 'get_campaign_by_id':
        await query.message.reply("Введите advertiser_id и campaign_id через пробел:")
        @dp.message()
        async def process_get_campaign_by_id(message: Message):
            args = message.text.strip().split()
            if len(args) != 2:
                await message.reply('Неверный формат. Введите advertiser_id и campaign_id через пробел.')
                return
            advertiser_id, campaign_id = args
            url = f"{API_BASE_URL}/advertisers/{advertiser_id}/campaigns/{campaign_id}"
            response = requests.get(url)
            if response.status_code == 200:
                campaign = response.json()
                reply_message = (
                    f"**Кампания ID**: `{campaign['campaign_id']}`\n"
                    f"**Название**: `{campaign['ad_title']}`\n"
                    f"**Текст**: `{campaign['ad_text']}`"
                )
                await message.reply(reply_message)
            else:
                await message.reply(f"Ошибка: {response.status_code}")
            dp.message_handlers.unregister(process_get_campaign_by_id)

    elif query.data == 'update_campaign':
        await query.message.reply("Отправьте данные для обновления кампании в формате JSON вместе с advertiser_id и campaign_id в URL:")
        @dp.message()
        async def process_update_campaign(message: Message):
            try:
                data = message.text.strip()
                advertiser_id = data.get('advertiser_id')
                campaign_id = data.get('campaign_id')
                url = f"{API_BASE_URL}/advertisers/{advertiser_id}/campaigns/{campaign_id}"
                campaign = requests.put(url, json=data)
                if campaign.status_code == 200:
                    await message.reply("Кампания успешно обновлена.")
                else:
                    await message.reply(f"Ошибка: {campaign.status_code}")
            except Exception as e:
                await message.reply(f"Ошибка: {e}")
            dp.message_handlers.unregister(process_update_campaign)

    elif query.data == 'delete_campaign':
        await query.message.reply("Введите advertiser_id и campaign_id через пробел:")
        @dp.message()
        async def process_delete_campaign(message: Message):
            args = message.text.strip().split()
            if len(args) != 2:
                await message.reply('Неверный формат. Введите advertiser_id и campaign_id через пробел.')
                return
            advertiser_id, campaign_id = args
            url = f"{API_BASE_URL}/advertisers/{advertiser_id}/campaigns/{campaign_id}"
            response = requests.delete(url)
            if response.status_code == 204:
                await message.reply("Кампания успешно удалена.")
            else:
                await message.reply(f"Ошибка: {response.status_code}")
            dp.message_handlers.unregister(process_delete_campaign)

    elif query.data == 'get_ad_for_client':
        await query.message.reply("Введите client_id:")
        @dp.message()
        async def process_get_ad_for_client(message: Message):
            client_id = message.text.strip()
            url = f"{API_BASE_URL}/ads"
            params = {'client_id': client_id}
            response = requests.get(url, params=params)
            if response.status_code == 200:
                ad = response.json()
                reply_message = (
                    f"**Реклама ID**: `{ad['ad_id']}`\n"
                    f"**Название**: `{ad['ad_title']}`\n"
                    f"**Текст**: `{ad['ad_text']}`\n"
                    f"**Рекламодатель ID**: `{ad['advertiser_id']}`"
                )
                await message.reply(reply_message)
            else:
                await message.reply(f"Ошибка: {response.status_code}")
            dp.message_handlers.unregister(process_get_ad_for_client)

    elif query.data == 'record_ad_click':
        await query.message.reply("Введите ad_id и client_id через пробел:")
        @dp.message()
        async def process_record_ad_click(message: Message):
            args = message.text.strip().split()
            if len(args) != 2:
                await message.reply('Неверный формат. Введите ad_id и client_id через пробел.')
                return
            ad_id, client_id = args
            url = f"{API_BASE_URL}/ads/{ad_id}/click"
            data = {'client_id': client_id}
            response = requests.post(url, json=data)
            if response.status_code == 204:
                await message.reply("Клик по объявлению успешно зафиксирован.")
            else:
                await message.reply(f"Ошибка: {response.status_code}")
            dp.message_handlers.unregister(process_record_ad_click)

    elif query.data == 'campaign_stats':
        await query.message.reply("Введите campaign_id:")
        @dp.message()
        async def process_campaign_stats(message: Message):
            campaign_id = message.text.strip()
            url = f"{API_BASE_URL}/stats/campaigns/{campaign_id}"
            response = requests.get(url)
            if response.status_code == 200:
                stats = response.json()
                reply_message = (
                    f"**Статистика кампании ID**: `{campaign_id}`\n"
                    f"**Показы**: `{stats['impressions_count']}`\n"
                    f"**Клики**: `{stats['clicks_count']}`\n"
                    f"**Конверсия**: `{stats['conversion']}`%\n"
                    f"**Затраты на показы**: `{stats['spent_impressions']}`\n"
                    f"**Затраты на клики**: `{stats['spent_clicks']}`\n"
                    f"**Общие затраты**: `{stats['spent_total']}`"
                )
                await message.reply(reply_message)
            else:
                await message.reply(f"Ошибка: {response.status_code}")
            dp.message_handlers.unregister(process_campaign_stats)

    elif query.data == 'advertiser_campaign_stats':
        await query.message.reply("Введите advertiser_id:")
        @dp.message()
        async def process_advertiser_campaign_stats(message: Message):
            advertiser_id = message.text.strip()
            url = f"{API_BASE_URL}/stats/advertisers/{advertiser_id}/campaigns"
            response = requests.get(url)
            if response.status_code == 200:
                stats = response.json()
                reply_message = (
                    f"**Статистика рекламодателя ID**: `{advertiser_id}`\n"
                    f"**Показы**: `{stats['impressions_count']}`\n"
                    f"**Клики**: `{stats['clicks_count']}`\n"
                    f"**Конверсия**: `{stats['conversion']}`%\n"
                    f"**Затраты на показы**: `{stats['spent_impressions']}`\n"
                    f"**Затраты на клики**: `{stats['spent_clicks']}`\n"
                    f"**Общие затраты**: `{stats['spent_total']}`"
                )
                await message.reply(reply_message)
            else:
                await message.reply(f"Ошибка: {response.status_code}")
            dp.message_handlers.unregister(process_advertiser_campaign_stats)

    elif query.data == 'campaign_daily_stats':
        await query.message.reply("Введите campaign_id:")
        @dp.message()
        async def process_campaign_daily_stats(message: Message):
            campaign_id = message.text.strip()
            url = f"{API_BASE_URL}/stats/campaigns/{campaign_id}/daily"
            response = requests.get(url)
            if response.status_code == 200:
                daily_stats = response.json()
                if not daily_stats:
                    await message.reply("Статистика не найдена.")
                    return
                reply_message = f"Ежедневная статистика кампании ID: `{campaign_id}`\n"
                for stats in daily_stats:
                    reply_message += (
                        f"**Дата**: `{stats['date']}`, "
                        f"**Показы**: `{stats['impressions_count']}`, "
                        f"**Клики**: `{stats['clicks_count']}`, "
                        f"**Конверсия**: `{stats['conversion']}`%, "
                        f"**Затраты**: `{stats['spent_total']}`\n"
                    )
                await message.reply(reply_message)
            else:
                await message.reply(f"Ошибка: {response.status_code}")
            dp.message_handlers.unregister(process_campaign_daily_stats)

    elif query.data == 'advertiser_daily_stats':
        await query.message.reply("Введите advertiser_id:")
        @dp.message()
        async def process_advertiser_daily_stats(message: Message):
            advertiser_id = message.text.strip()
            url = f"{API_BASE_URL}/stats/advertisers/{advertiser_id}/campaigns/daily"
            response = requests.get(url)
            if response.status_code == 200:
                daily_stats = response.json()
                if not daily_stats:
                    await message.reply("Статистика не найдена.")
                    return
                reply_message = f"Ежедневная статистика рекламодателя ID: `{advertiser_id}`\n"
                for stats in daily_stats:
                    reply_message += (
                        f"**Дата**: `{stats['date']}`, "
                        f"**Показы**: `{stats['impressions_count']}`, "
                        f"**Клики**: `{stats['clicks_count']}`, "
                        f"**Конверсия**: `{stats['conversion']}`%, "
                        f"**Затраты**: `{stats['spent_total']}`\n"
                    )
                await message.reply(reply_message)
            else:
                await message.reply(f"Ошибка: {response.status_code}")
            dp.message_handlers.unregister(process_advertiser_daily_stats)

    elif query.data == 'advance_time':
        await query.message.reply("Введите новую текущую дату (целое число):")
        @dp.message()
        async def process_advance_time(message: Message):
            current_date = message.text.strip()
            try:
                data = {'current_date': int(current_date)}
                url = f"{API_BASE_URL}/time/advance"
                response = requests.post(url, json=data)
                if response.status_code == 200:
                    new_date = response.json().get('current_date')
                    await message.reply(f"Текущая дата установлена: {new_date}")
                else:
                    await message.reply(f"Ошибка: {response.status_code}")
            except ValueError:
                await message.reply("Пожалуйста, введите корректное целое число.")
            dp.message_handlers.unregister(process_advance_time)

# Обработчик команды /get_client
@dp.message(Command(commands=['get_client']))
async def get_client(message: Message):
    args = message.text.split()
    if len(args) != 2:
        await message.reply('Использование: /get_client <client_id>')
        return
    client_id = args[1]
    url = f"{API_BASE_URL}/clients/{client_id}"
    response = requests.get(url)
    if response.status_code == 200:
        client = response.json()
        reply_message = (
            f"**Клиент ID**: `{client['client_id']}`\n"
            f"**Логин**: `{client['login']}`\n"
            f"**Возраст**: `{client['age']}`\n"
            f"**Локация**: `{client['location']}`\n"
            f"**Пол**: `{client['gender']}`"
        )
        await message.reply(reply_message)
    else:
        await message.reply(f"Ошибка: {response.status_code}")

# Обработчик команды /get_advertiser
@dp.message(Command(commands=['get_advertiser']))
async def get_advertiser(message: Message):
    args = message.text.split()
    if len(args) != 2:
        await message.reply('Использование: /get_advertiser <advertiser_id>')
        return
    advertiser_id = args[1]
    url = f"{API_BASE_URL}/advertisers/{advertiser_id}"
    response = requests.get(url)
    if response.status_code == 200:
        advertiser = response.json()
        reply_message = (
            f"**Рекламодатель ID**: `{advertiser['advertiser_id']}`\n"
            f"**Название**: `{advertiser['name']}`"
        )
        await message.reply(reply_message)
    else:
        await message.reply(f"Ошибка: {response.status_code}")

# Обработчик команды /get_campaigns
@dp.message(Command(commands=['get_campaigns']))
async def get_campaigns(message: Message):
    args = message.text.split()
    if len(args) < 2:
        await message.reply('Использование: /get_campaigns <advertiser_id> [page] [size]')
        return
    advertiser_id = args[1]
    params = {}
    if len(args) > 2:
        params['page'] = args[2]
    if len(args) > 3:
        params['size'] = args[3]
    url = f"{API_BASE_URL}/advertisers/{advertiser_id}/campaigns"
    response = requests.get(url, params=params)
    if response.status_code == 200:
        campaigns = response.json()
        if not campaigns:
            await message.reply("Кампании не найдены.")
            return
        reply_message = "Список кампаний:\n"
        for campaign in campaigns:
            reply_message += (
                f"-**Кампания ID**: `{campaign['campaign_id']}`, "
                f"**Название**: `{campaign['ad_title']}`\n"
            )
        await message.reply(reply_message)
    else:
        await message.reply(f"Ошибка: {response.status_code}")

# Обработчик команды /get_campaign
@dp.message(Command(commands=['get_campaign']))
async def get_campaign(message: Message):
    args = message.text.split()
    if len(args) != 3:
        await message.reply('Использование: /get_campaign <advertiser_id> <campaign_id>')
        return
    advertiser_id = args[1]
    campaign_id = args[2]
    url = f"{API_BASE_URL}/advertisers/{advertiser_id}/campaigns/{campaign_id}"
    response = requests.get(url)
    if response.status_code == 200:
        campaign = response.json()
        reply_message = (
            f"**Кампания ID**: `{campaign['campaign_id']}`\n"
            f"**Название**: `{campaign['ad_title']}`\n"
            f"**Текст**: `{campaign['ad_text']}`"
        )
        await message.reply(reply_message)
    else:
        await message.reply(f"Ошибка: {response.status_code}")

# Обработчик команды /get_ad
@dp.message(Command(commands=['get_ad']))
async def get_ad(message: Message):
    args = message.text.split()
    if len(args) != 2:
        await message.reply('Использование: /get_ad <client_id>')
        return
    client_id = args[1]
    url = f"{API_BASE_URL}/ads"
    params = {'client_id': client_id}
    response = requests.get(url, params=params)
    if response.status_code == 200:
        ad = response.json()
        reply_message = (
            f"**Реклама ID**: `{ad['ad_id']}`\n"
            f"**Название**: `{ad['ad_title']}`\n"
            f"**Текст**: `{ad['ad_text']}`\n"
            f"**Рекламодатель ID**: `{ad['advertiser_id']}`"
        )
        await message.reply(reply_message)
    else:
        await message.reply(f"Ошибка: {response.status_code}")

# Обработчик команды /get_campaign_stats
@dp.message(Command(commands=['get_campaign_stats']))
async def get_campaign_stats(message: Message):
    args = message.text.split()
    if len(args) != 2:
        await message.reply('Использование: /get_campaign_stats <campaign_id>')
        return
    campaign_id = args[1]
    url = f"{API_BASE_URL}/stats/campaigns/{campaign_id}"
    response = requests.get(url)
    if response.status_code == 200:
        stats = response.json()
        reply_message = (
            f"**Статистика кампании ID**: `{campaign_id}`\n"
            f"**Показы**: `{stats['impressions_count']}`\n"
            f"**Клики**: `{stats['clicks_count']}`\n"
            f"**Конверсия**: `{stats['conversion']}`%\n"
            f"**Затраты на показы**: `{stats['spent_impressions']}`\n"
            f"**Затраты на клики**: `{stats['spent_clicks']}`\n"
            f"**Общие затраты**: `{stats['spent_total']}`"
        )
        await message.reply(reply_message)
    else:
        await message.reply(f"Ошибка: {response.status_code}")

# Обработчик команды /get_advertiser_stats
@dp.message(Command(commands=['get_advertiser_stats']))
async def get_advertiser_stats(message: Message):
    args = message.text.split()
    if len(args) != 2:
        await message.reply('Использование: /get_advertiser_stats <advertiser_id>')
        return
    advertiser_id = args[1]
    url = f"{API_BASE_URL}/stats/advertisers/{advertiser_id}/campaigns"
    response = requests.get(url)
    if response.status_code == 200:
        stats = response.json()
        reply_message = (
            f"**Статистика рекламодателя ID**: `{advertiser_id}`\n"
            f"**Показы**: `{stats['impressions_count']}`\n"
            f"**Клики**: `{stats['clicks_count']}`\n"
            f"**Конверсия**: `{stats['conversion']}`%\n"
            f"**Затраты на показы**: `{stats['spent_impressions']}`\n"
            f"**Затраты на клики**: `{stats['spent_clicks']}`\n"
            f"**Общие затраты**: `{stats['spent_total']}`"
        )
        await message.reply(reply_message)
    else:
        await message.reply(f"Ошибка: {response.status_code}")

# Обработчик команды /get_campaign_daily_stats
@dp.message(Command(commands=['get_campaign_daily_stats']))
async def get_campaign_daily_stats(message: Message):
    args = message.text.split()
    if len(args) != 2:
        await message.reply('Использование: /get_campaign_daily_stats <campaign_id>')
        return
    campaign_id = args[1]
    url = f"{API_BASE_URL}/stats/campaigns/{campaign_id}/daily"
    response = requests.get(url)
    if response.status_code == 200:
        daily_stats = response.json()
        if not daily_stats:
            await message.reply("Статистика не найдена.")
            return
        reply_message = f"Ежедневная статистика кампании ID: `{campaign_id}`\n"
        for stats in daily_stats:
            reply_message += (
                f"**Дата**: `{stats['date']}`, "
                f"**Показы**: `{stats['impressions_count']}`, "
                f"**Клики**: `{stats['clicks_count']}`, "
                f"**Конверсия**: `{stats['conversion']}`%, "
                f"**Затраты**: `{stats['spent_total']}`\n"
            )
        await message.reply(reply_message)
    else:
        await message.reply(f"Ошибка: {response.status_code}")

# Обработчик команды /get_advertiser_daily_stats
@dp.message(Command(commands=['get_advertiser_daily_stats']))
async def get_advertiser_daily_stats(message: Message):
    args = message.text.split()
    if len(args) != 2:
        await message.reply('Использование: /get_advertiser_daily_stats <advertiser_id>')
        return
    advertiser_id = args[1]
    url = f"{API_BASE_URL}/stats/advertisers/{advertiser_id}/campaigns/daily"
    response = requests.get(url)
    if response.status_code == 200:
        daily_stats = response.json()
        if not daily_stats:
            await message.reply("Статистика не найдена.")
            return
        reply_message = f"Ежедневная статистика рекламодателя ID: `{advertiser_id}`\n"
        for stats in daily_stats:
            reply_message += (
                f"**Дата**: `{stats['date']}`, "
                f"**Показы**: `{stats['impressions_count']}`, "
                f"**Клики**: `{stats['clicks_count']}`, "
                f"**Конверсия**: `{stats['conversion']}`%, "
                f"**Затраты**: `{stats['spent_total']}`\n"
            )
        await message.reply(reply_message)
    else:
        await message.reply(f"Ошибка: {response.status_code}")

# Главная функция
async def main():
    # Запуск бота
    try:
        await dp.start_polling(bot)
    finally:
        await bot.session.close()

if __name__== '__main__':
    asyncio.run(main())
