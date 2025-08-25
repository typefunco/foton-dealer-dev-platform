# Dealer Dev Platform Frontend

Modern React application for managing dealer network using TypeScript, Tailwind CSS and Framer Motion.

## 🚀 Technologies

- **React 18.2.0** - Main library
- **TypeScript 5.1.6** - Type safety
- **Vite 4.4.5** - Build tool and dev server
- **Tailwind CSS 3.3.3** - CSS framework
- **Framer Motion 10.16.4** - Animations
- **React Router 6.8.1** - Routing
- **Heroicons** - Icons

## 📦 Install Dependencies

```bash
yarn install
```

## 🏃‍♂️ Development Mode

```bash
yarn dev
```

The application will be available at: http://localhost:3000

## 🏗️ Production Build

```bash
yarn build
```

## 📱 Preview Build

```bash
yarn preview
```

## 🔍 Code Check

```bash
yarn lint
```

## 📁 Project Structure

```
src/
├── components/     # Reusable components
├── pages/         # Application pages
├── hooks/         # Custom hooks
├── utils/         # Utilities
├── types/         # TypeScript types
├── assets/        # Static resources
├── App.tsx        # Main component
├── main.tsx       # Entry point
└── index.css      # Global styles
```

## 🎨 Особенности дизайна

- **Responsive дизайн** - Адаптация под все устройства
- **Плавные анимации** - Использование Framer Motion
- **Современный UI** - Tailwind CSS компоненты
- **Иконки** - Heroicons

## 🚀 Основные страницы

1. **Главная** (`/`) - Лендинг с описанием возможностей

## 🎯 Компоненты

- **App** - Главный компонент с анимациями
- **Button** - Кнопки с анимациями

## 🔧 Настройка

### Tailwind CSS
Конфигурация находится в `tailwind.config.js` с базовыми настройками.

### TypeScript
Строгая типизация с настройками в `tsconfig.json`.

### Vite
Быстрая сборка и hot reload в `vite.config.ts`.

## 📱 Поддерживаемые браузеры

- Chrome (последние 2 версии)
- Firefox (последние 2 версии)
- Safari (последние 2 версии)
- Edge (последние 2 версии)

## 🤝 Разработка

### Добавление новых страниц
1. Создайте компонент в `src/pages/`
2. Добавьте роут в `src/App.tsx`
3. Добавьте ссылку в навигацию

### Добавление новых компонентов
1. Создайте компонент в `src/components/`
2. Используйте TypeScript интерфейсы
3. Добавьте анимации с Framer Motion

### Стилизация
- Используйте Tailwind CSS классы
- Создавайте кастомные компоненты в `src/index.css`
- Следуйте дизайн-системе проекта

## 🚨 Известные проблемы

- Нет проблем на данный момент

## 📄 Лицензия

MIT License
