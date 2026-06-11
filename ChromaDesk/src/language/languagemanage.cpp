#include "languagemanage.h"

#ifdef Q_OS_WIN32
#include <Windows.h>
#endif

#include <QGuiApplication>
#include <QTranslator>

#include "core/localconfigcenter.h"

static QString kLanguagePath = ":/i18n/";

LanguageManage::LanguageManage(QObject* parent)
    : QObject(parent)
{
}

LanguageManage::~LanguageManage()
{
}

void LanguageManage::init()
{
    installLanguage();
}

QStringList LanguageManage::getSupportLanguages()
{
    QStringList languages;
    languages << "Auto" << "zh_CN" << "en_US";
    return languages;
}

void LanguageManage::setCurrentLanguage(const QString& language)
{
    core::LocalConfigCenter::instance().setLanguage(language);
    installLanguage();
}

QString LanguageManage::getCurrentLanguage()
{
    return core::LocalConfigCenter::instance().language();
}

QString LanguageManage::getCurrentRealLanguage()
{
    QString language = getCurrentLanguage();
    if (language == "Auto") {
        QLocale::Language systemLanguage = getSystemLanguage();
        if (systemLanguage == QLocale::Chinese) {
            language = "zh_CN";
        } else if (systemLanguage == QLocale::English) {
            language = "en_US";
        }
    }
    return language;
}

QString LanguageManage::getLanguageName(const QString& language)
{
    if (language == "zh_CN") {
        return u8"简体中文";
    }
    if (language == "en_US") {
        return "English";
    }
    return tr("Auto");
}

void LanguageManage::installLanguage()
{
    auto configLanguage = core::LocalConfigCenter::instance().language();
    QLocale::Language language = getSystemLanguage();
    if (configLanguage == "zh_CN") {
        language = QLocale::Chinese;
    }
    if (configLanguage == "en_US") {
        language = QLocale::English;
    }

    installTranslator(language);
}

QLocale::Language LanguageManage::getSystemLanguage()
{
    QLocale locale;
    QLocale::Language language = locale.language();

#ifdef Q_OS_WIN32
    LANGID languageId = ::GetUserDefaultUILanguage() & 0x03FF;
    if (LANG_CHINESE_SIMPLIFIED == languageId || LANG_CHINESE_TRADITIONAL == languageId) {
        language = QLocale::Chinese;
    }
#endif

    return language;
}

void LanguageManage::installTranslator(const QLocale::Language& language)
{
    static QTranslator zhCnTranslator;
    static QTranslator enUsTranslator;
    QTranslator* translator = nullptr;
    if (language == QLocale::Chinese) {
        if (zhCnTranslator.isEmpty()) {
            (void)zhCnTranslator.load(kLanguagePath + "zh_CN.qm");
        }
        translator = &zhCnTranslator;
    } else {
        if (enUsTranslator.isEmpty()) {
            (void)enUsTranslator.load(kLanguagePath + "en_US.qm");
        }
        translator = &enUsTranslator;
    }

    qApp->installTranslator(translator);
}
