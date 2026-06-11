// Fluent Design Badge Component
import QtQuick

Rectangle {
    id: control
    
    // ============ Badge Types ============
    
    enum Type {
        Primary,
        Success,
        Warning,
        Error,
        Info
    }
    
    // ============ Custom Properties ============
    
    property int badgeType: QDBadge.Type.Primary
    property string text: ""
    property int count: -1  // -1 means show text instead of count
    property int maxCount: 99
    property bool dot: false  // Show as dot only
    
    // ============ Size & Style ============
    
    implicitWidth: dot ? 8 : Math.max(20, badgeContent.width + Theme.spacingSmall * 2)
    implicitHeight: dot ? 8 : 20
    radius: height / 2
    
    color: {
        switch(badgeType) {
            case QDBadge.Type.Primary: return Theme.primary
            case QDBadge.Type.Success: return Theme.success
            case QDBadge.Type.Warning: return Theme.warning
            case QDBadge.Type.Error: return Theme.error
            case QDBadge.Type.Info: return Theme.info
            default: return Theme.primary
        }
    }
    
    // ============ Content ============
    
    Text {
        id: badgeContent
        visible: !dot
        anchors.centerIn: parent
        text: {
            if (control.count >= 0) {
                return control.count > control.maxCount ? 
                       control.maxCount + "+" : control.count.toString()
            }
            return control.text
        }
        font.family: Theme.fontFamily
        font.pixelSize: Theme.fontSizeSmall
        font.weight: Font.DemiBold
        color: Theme.textOnPrimary
    }
}
