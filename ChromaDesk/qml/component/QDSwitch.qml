// Fluent Design Switch Component
import QtQuick
import QtQuick.Controls as Controls

Controls.Switch {
    id: control
    
    // ============ Size & Style ============
    
    implicitWidth: indicator.width + leftPadding + rightPadding + 
                   (text !== "" ? contentItem.implicitWidth + spacing : 0)
    implicitHeight: Math.max(indicator.height, contentItem.implicitHeight) + topPadding + bottomPadding
    
    spacing: Theme.spacingMedium
    padding: Theme.spacingSmall
    
    // ============ Indicator (Switch track and handle) ============
    
    indicator: Rectangle {
        id: track
        implicitWidth: 44
        implicitHeight: 22
        x: control.leftPadding
        y: parent.height / 2 - height / 2
        radius: height / 2
        
        color: {
            if (!control.enabled) {
                return Theme.surfaceVariant
            }
            if (control.checked) {
                return control.down ? Theme.primaryPressed :
                       control.hovered ? Theme.primaryHover : Theme.primary
            }
            return control.hovered ? Theme.surfaceHover : Theme.surface
        }
        
        border.width: control.checked ? 0 : Theme.borderWidthMedium
        border.color: control.enabled ? Theme.border : Theme.border
        
        Behavior on color {
            ColorAnimation {
                duration: Theme.animationDurationFast
                easing.type: Theme.animationEasingType
            }
        }
        
        // Handle (thumb)
        Rectangle {
            id: handle
            width: 16
            height: 16
            radius: 8
            x: control.checked ? parent.width - width - 3 : 3
            anchors.verticalCenter: parent.verticalCenter
            
            color: {
                if (!control.enabled) {
                    return Theme.textDisabled
                }
                return Theme.textOnPrimary
            }
            
            Behavior on x {
                NumberAnimation {
                    duration: Theme.animationDurationFast
                    easing.type: Easing.OutCubic
                }
            }
            
            // Inner shadow effect
            Rectangle {
                anchors.fill: parent
                radius: parent.radius
                color: "transparent"
                border.width: 1
                border.color: Qt.rgba(0, 0, 0, 0.1)
            }
        }
    }
    
    // ============ Content (Text label) ============
    
    contentItem: Text {
        leftPadding: control.indicator.width + control.spacing
        text: control.text
        font.family: Theme.fontFamily
        font.pixelSize: Theme.fontSizeMedium
        color: control.enabled ? Theme.text : Theme.textDisabled
        verticalAlignment: Text.AlignVCenter
    }
    
    // Hover cursor
    HoverHandler {
        cursorShape: control.enabled ? Qt.PointingHandCursor : Qt.ForbiddenCursor
    }
}
