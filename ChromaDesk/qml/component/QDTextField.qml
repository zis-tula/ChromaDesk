// Fluent Design TextField Component
import QtQuick
import QtQuick.Controls as Controls

Controls.TextField {
    id: control
    
    // ============ Custom Properties ============
    
    property string label: ""
    property bool showLabel: label !== ""
    property bool error: false
    property string errorText: ""
    property string prefixIcon: ""
    property string suffixIcon: ""
    
    // ============ Size & Style ============
    
    implicitWidth: 200
    implicitHeight: Theme.buttonHeightMedium
    
    leftPadding: prefixIcon !== "" ? Theme.spacingXLarge + Theme.iconSizeMedium : Theme.spacingMedium
    rightPadding: suffixIcon !== "" ? Theme.spacingXLarge + Theme.iconSizeMedium : Theme.spacingMedium
    topPadding: Theme.spacingSmall
    bottomPadding: Theme.spacingSmall
    
    color: control.enabled ? Theme.text : Theme.textDisabled
    selectionColor: Theme.primary
    selectedTextColor: Theme.textOnPrimary
    placeholderTextColor: Theme.textSecondary
    
    font.family: Theme.fontFamily
    font.pixelSize: Theme.fontSizeMedium
    
    verticalAlignment: TextInput.AlignVCenter
    
    // ============ Background ============
    
    background: Rectangle {
        id: textFieldBackground
        radius: Theme.radiusMedium
        color: control.enabled ? Theme.surface : Theme.surfaceVariant
        border.width: Theme.borderWidthMedium
        border.color: {
            if (control.error) {
                return Theme.error
            }
            if (control.activeFocus) {
                return Theme.borderFocus
            }
            if (control.hovered) {
                return Theme.borderHover
            }
            return Theme.border
        }
        
        Behavior on border.color {
            ColorAnimation {
                duration: Theme.animationDurationFast
                easing.type: Theme.animationEasingType
            }
        }
        
        Behavior on color {
            ColorAnimation {
                duration: Theme.animationDurationFast
                easing.type: Theme.animationEasingType
            }
        }
        
        // Prefix icon
        Text {
            visible: prefixIcon !== ""
            text: prefixIcon
            font.family: "Segoe Fluent Icons"
            font.pixelSize: Theme.iconSizeMedium
            color: control.activeFocus ? Theme.primary : Theme.textSecondary
            anchors.left: parent.left
            anchors.leftMargin: Theme.spacingMedium
            anchors.verticalCenter: parent.verticalCenter
            
            Behavior on color {
                ColorAnimation { duration: Theme.animationDurationFast }
            }
        }
        
        // Suffix icon
        Text {
            visible: suffixIcon !== ""
            text: suffixIcon
            font.family: "Segoe Fluent Icons"
            font.pixelSize: Theme.iconSizeMedium
            color: control.activeFocus ? Theme.primary : Theme.textSecondary
            anchors.right: parent.right
            anchors.rightMargin: Theme.spacingMedium
            anchors.verticalCenter: parent.verticalCenter
            
            Behavior on color {
                ColorAnimation { duration: Theme.animationDurationFast }
            }
        }
        
        // Focus underline animation
        Rectangle {
            id: focusLine
            anchors.bottom: parent.bottom
            anchors.bottomMargin: Theme.borderWidthMedium
            anchors.horizontalCenter: parent.horizontalCenter
            height: 2
            // 限制宽度不超出圆角区域
            width: control.activeFocus ? parent.width - textFieldBackground.radius * 2 : 0
            color: Theme.primary
            radius: 1
            
            Behavior on width {
                NumberAnimation {
                    duration: Theme.animationDurationMedium
                    easing.type: Easing.OutCubic
                }
            }
        }
    }
    
    // Hover cursor
    HoverHandler {
        cursorShape: control.enabled ? Qt.IBeamCursor : Qt.ForbiddenCursor
    }
}
