# Fluent Design 组件库

ChromaDesk 的 Modern Fluent Design 风格 UI 组件库，支持多主题切换。

## 主题系统

### 内置主题

| 主题名称 | 类型 | 说明 |
|---------|------|------|
| **Fluent Dark** | 深色 | 默认深色主题，蓝色强调（微软标准） |
| **Fluent Light** | 浅色 | 浅色主题，蓝色强调 |
| **Nord** | 深色 | Nord 配色方案，青色强调 |
| **Dracula** | 深色 | Dracula 配色方案，紫色强调 |
| **Monokai** | 深色 | Monokai 配色方案，绿色强调 |
| **Solarized Light** | 浅色 | Solarized 配色方案，橙色强调 |

### 切换主题

```qml
// 在 QML 中切换主题
Theme.currentTheme = Theme.ThemeType.DraculaDark

// 可用的主题类型
Theme.ThemeType.FluentDark        // 默认
Theme.ThemeType.FluentLight
Theme.ThemeType.NordDark
Theme.ThemeType.DraculaDark
Theme.ThemeType.MonokaiDark
Theme.ThemeType.SolarizedLight
```

### 获取当前主题名称

```qml
Text {
    text: "当前主题: " + Theme.currentThemeName
}
```

### 主题切换动画

所有颜色属性都有 250ms 的平滑过渡动画，主题切换时整个界面颜色会平滑变化。

---

## 组件列表

### 基础组件

| 组件 | 文件 | 说明 |
|------|------|------|
| **Theme** | Theme.qml | 主题管理器（Singleton，支持6种主题） |
| **FluentIconGlyph** | FluentIconGlyph.qml | 图标定义（Singleton） |

### 按钮与输入

| 组件 | 文件 | 说明 |
|------|------|------|
| **QDButton** | QDButton.qml | 按钮（Primary/Secondary/Danger/Success/Ghost） |
| **QDTextField** | QDTextField.qml | 文本输入框 |
| **QDCheckBox** | QDCheckBox.qml | 复选框 |
| **QDSwitch** | QDSwitch.qml | 开关 |
| **QDSlider** | QDSlider.qml | 滑块 |
| **QDComboBox** | QDComboBox.qml | 下拉选择框 |

### 容器与布局

| 组件 | 文件 | 说明 |
|------|------|------|
| **QDCard** | QDCard.qml | 卡片容器 |
| **QDDialog** | QDDialog.qml | 对话框 |
| **QDDivider** | QDDivider.qml | 分割线 |

### 反馈与提示

| 组件 | 文件 | 说明 |
|------|------|------|
| **QDToast** | QDToast.qml | 消息提示 |
| **QDProgressBar** | QDProgressBar.qml | 进度条 |
| **QDBadge** | QDBadge.qml | 徽章 |
| **QDToolTip** | QDToolTip.qml | 工具提示 |

---

## 使用方法

### 1. 导入组件

```qml
import "component"
```

### 2. 使用主题

所有组件都会自动使用 Theme 中定义的颜色，支持主题切换：

```qml
Rectangle {
    color: Theme.background  // 自动跟随主题变化
}

Text {
    color: Theme.text
    font.family: Theme.fontFamily
}

QDButton {
    // 按钮会自动使用当前主题的颜色
    text: "按钮"
    buttonType: QDButton.Type.Primary
}
```

### 3. 创建主题切换器

```qml
Row {
    spacing: 8
    
    Repeater {
        model: [
            { name: "Dark", type: Theme.ThemeType.FluentDark },
            { name: "Light", type: Theme.ThemeType.FluentLight },
            { name: "Nord", type: Theme.ThemeType.NordDark },
            { name: "Dracula", type: Theme.ThemeType.DraculaDark },
            { name: "Monokai", type: Theme.ThemeType.MonokaiDark },
            { name: "Solarized", type: Theme.ThemeType.SolarizedLight }
        ]
        
        QDButton {
            text: modelData.name
            buttonType: Theme.currentTheme === modelData.type ? 
                       QDButton.Type.Primary : QDButton.Type.Secondary
            onClicked: Theme.currentTheme = modelData.type
        }
    }
}
```

---

## 组件详细 API

### QDButton

| 属性 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `text` | string | "" | 按钮文本 |
| `buttonType` | enum | Primary | 按钮类型 |
| `iconText` | string | "" | 图标（使用 FluentIconGlyph） |
| `loading` | bool | false | 加载状态 |
| `enabled` | bool | true | 是否可用 |

**buttonType 枚举值：**
- `QDButton.Type.Primary` - 主要按钮
- `QDButton.Type.Secondary` - 次要按钮
- `QDButton.Type.Danger` - 危险按钮
- `QDButton.Type.Success` - 成功按钮
- `QDButton.Type.Ghost` - 幽灵按钮

### QDTextField

| 属性 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `placeholderText` | string | "" | 占位文本 |
| `prefixIcon` | string | "" | 前缀图标 |
| `suffixIcon` | string | "" | 后缀图标 |
| `error` | bool | false | 错误状态 |
| `errorText` | string | "" | 错误提示文本 |

### QDCheckBox

| 属性 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `text` | string | "" | 标签文本 |
| `checked` | bool | false | 选中状态 |
| `checkState` | enum | Qt.Unchecked | 三态状态 |

### QDSwitch

| 属性 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `text` | string | "" | 标签文本 |
| `checked` | bool | false | 开关状态 |

### QDSlider

| 属性 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `from` | real | 0.0 | 最小值 |
| `to` | real | 1.0 | 最大值 |
| `value` | real | 0.0 | 当前值 |
| `showValue` | bool | false | 拖动时显示值 |
| `decimals` | int | 0 | 小数位数 |

### QDComboBox

| 属性 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `model` | var | [] | 数据模型 |
| `currentIndex` | int | 0 | 当前选中索引 |
| `displayText` | string | - | 显示文本 |

### QDCard

| 属性 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `elevation` | int | 1 | 阴影级别（1-3） |
| `hoverable` | bool | false | 是否可悬停 |
| `clickable` | bool | false | 是否可点击 |

**信号：**
- `clicked()` - 点击时触发

### QDToast

| 属性 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `toastType` | enum | Info | 提示类型 |
| `message` | string | "" | 消息内容 |
| `duration` | int | 3000 | 显示时长（毫秒） |
| `autoHide` | bool | true | 自动隐藏 |

**方法：**
- `show(message, type, duration)` - 显示提示
- `hide()` - 隐藏提示

**toastType 枚举值：**
- `QDToast.Type.Success` - 成功（绿色）
- `QDToast.Type.Error` - 错误（红色）
- `QDToast.Type.Warning` - 警告（橙色）
- `QDToast.Type.Info` - 信息（蓝色）

### QDDialog

| 属性 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `title` | string | "Dialog" | 对话框标题 |
| `showing` | bool | false | 显示状态 |
| `dialogWidth` | int | 400 | 宽度 |
| `dialogHeight` | int | 300 | 高度 |
| `closeOnOverlay` | bool | true | 点击遮罩关闭 |
| `showCloseButton` | bool | true | 显示关闭按钮 |
| `footer` | list | [] | 底部按钮列表 |

**方法：**
- `show()` - 显示对话框
- `hide()` - 隐藏对话框
- `accept()` - 确认并关闭
- `reject()` - 取消并关闭

**信号：**
- `accepted()` - 确认时触发
- `rejected()` - 取消时触发
- `closed()` - 关闭时触发

### QDProgressBar

| 属性 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `value` | real | 0.0 | 当前进度（0-1） |
| `indeterminate` | bool | false | 不确定模式 |
| `showValue` | bool | false | 显示百分比 |
| `progressColor` | color | primary | 进度条颜色 |

### QDBadge

| 属性 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `badgeType` | enum | Primary | 徽章类型 |
| `text` | string | "" | 文本内容 |
| `count` | int | -1 | 数字（-1时显示text） |
| `maxCount` | int | 99 | 最大数字 |
| `dot` | bool | false | 圆点模式 |

**badgeType 枚举值：**
- `QDBadge.Type.Primary` - 主色
- `QDBadge.Type.Success` - 成功
- `QDBadge.Type.Warning` - 警告
- `QDBadge.Type.Error` - 错误
- `QDBadge.Type.Info` - 信息

### QDDivider

| 属性 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `vertical` | bool | false | 垂直方向 |
| `label` | string | "" | 中间标签文本 |

### QDToolTip

| 属性 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `text` | string | "" | 提示文本 |
| `tipDelay` | int | 500 | 延迟显示（毫秒） |

---

## 主题配置

### 主题 API

```qml
// 切换主题
Theme.currentTheme = Theme.ThemeType.NordDark

// 获取当前主题名称
Theme.currentThemeName  // "Nord Dark"

// 访问主题颜色（所有颜色都会随主题变化）
Theme.background
Theme.surface
Theme.primary
Theme.accent
Theme.text
Theme.textSecondary
Theme.success
Theme.error
Theme.warning
Theme.info
```

### 主题颜色列表

每个主题包含以下颜色：

```qml
background          // 背景色
surface             // 表面色（卡片、输入框等）
surfaceVariant      // 次级表面色
surfaceHover        // 表面悬停色
primary             // 主色
primaryHover        // 主色悬停
primaryPressed      // 主色按下
primaryDisabled     // 主色禁用
accent              // 强调色
accentLight         // 浅强调色
border              // 边框色
borderHover         // 边框悬停
borderFocus         // 边框聚焦
text                // 主文本色
textSecondary       // 次要文本色
textDisabled        // 禁用文本色
textOnPrimary       // 主色上的文本色
success             // 成功色
successHover        // 成功悬停
error               // 错误色
errorHover          // 错误悬停
warning             // 警告色
warningHover        // 警告悬停
info                // 信息色
infoHover           // 信息悬停
shadowLight         // 浅阴影
shadowMedium        // 中阴影
shadowDark          // 深阴影
overlay             // 遮罩层
```

### 尺寸（不随主题变化）

```qml
Theme.radiusSmall     // 4px
Theme.radiusMedium    // 8px
Theme.radiusLarge     // 12px
Theme.spacingSmall    // 8px
Theme.spacingMedium   // 12px
Theme.spacingLarge    // 16px
Theme.spacingXLarge   // 24px
Theme.spacingXXLarge  // 32px
```

### 动画（不随主题变化）

```qml
Theme.animationDurationFast    // 150ms
Theme.animationDurationMedium  // 250ms
Theme.animationDurationSlow    // 350ms
```

---

## 自定义主题

如果需要添加自定义主题，可以修改 `Theme.qml`：

```qml
// 1. 在 ThemeType 枚举中添加新主题
enum ThemeType {
    FluentDark,
    FluentLight,
    // ... 其他主题
    CustomTheme  // 新增
}

// 2. 添加主题配置对象
readonly property var customTheme: ({
    name: "Custom Theme",
    background: "#...",
    surface: "#...",
    // ... 其他颜色
})

// 3. 在 _currentThemeConfig 中添加 case
readonly property var _currentThemeConfig: {
    switch(currentTheme) {
        // ... 其他 case
        case Theme.ThemeType.CustomTheme: return customTheme
        default: return fluentDark
    }
}
```

---

## 文件结构

```
qml/component/
├── Theme.qml            # 主题管理器（支持6种主题）
├── FluentIconGlyph.qml  # 图标定义
├── QDButton.qml         # 按钮
├── QDTextField.qml      # 输入框
├── QDCheckBox.qml       # 复选框
├── QDSwitch.qml         # 开关
├── QDSlider.qml         # 滑块
├── QDComboBox.qml       # 下拉框
├── QDCard.qml           # 卡片
├── QDDialog.qml         # 对话框
├── QDToast.qml          # 消息提示
├── QDProgressBar.qml    # 进度条
├── QDBadge.qml          # 徽章
├── QDDivider.qml        # 分割线
├── QDToolTip.qml        # 工具提示
├── qmldir               # 模块定义
└── README.md            # 本文档
```
