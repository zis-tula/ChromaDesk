#pragma once

#include <QLocale>
#include <QObject>
#include <QStringList>

#include "base/singleton.h"

class LanguageManage : public QObject, public base::Singleton<LanguageManage> {
    Q_OBJECT
public:
    LanguageManage(QObject* parent = nullptr);
    ~LanguageManage();

    void init();
    Q_INVOKABLE QStringList getSupportLanguages();
    Q_INVOKABLE void setCurrentLanguage(const QString& language);
    Q_INVOKABLE QString getCurrentLanguage();
    // 和getCurrentLanguage的区别是"Auto"会转换为实际的语言
    Q_INVOKABLE QString getCurrentRealLanguage();
    Q_INVOKABLE QString getLanguageName(const QString& language);

private:
    void installLanguage();
    QLocale::Language getSystemLanguage();
    void installTranslator(const QLocale::Language& language);
};
