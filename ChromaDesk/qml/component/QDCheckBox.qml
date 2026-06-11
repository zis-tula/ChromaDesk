// Fluent Design CheckBox Component
import QtQuick
import QtQuick.Controls as Controls

Controls.CheckBox {
    id: control
    
    // ============ Size & Style ============
    
    implicitWidth: Math.max(indicator.width + leftPadding + rightPadding,
                           contentItem.implicitWidth + leftPadding + rightPadding)
    implicitHeight: Math.max(indicator.height + topPadding + bottomPadding,
                            contentItem.implicitHeight + topPadding + bottomPadding)
    
    spacing: Theme.spacingSmall
    padding: Theme.spacingSmall
    
    // ============ Indicator (Checkbox box) ============
    
    indicator: Rectangle {
        id: checkboxRect
        implicitWidth: 20
        implicitHeight: 20
        x: control.leftPadding
        y: parent.height / 2 - height / 2
        radius: Theme.radiusSmall
        
        border.width: Theme.borderWidthMedium
        border.color: {
            if (!control.enabled) {
                return Theme.border
            }
            if (control.checked) {
                return Theme.primary
            }
            if (control.hovered) {
                return Theme.borderHover
            }
            return Theme.border
        }
        
        color: {
            if (!control.enabled) {
                return Theme.surfaceVariant
            }
            if (control.checked) {
                return control.down ? Theme.primaryPressed :
                       control.hovered ? Theme.primaryHover : Theme.primary
            }
            if (control.down) {
                return Theme.surfaceVariant
            }
            if (control.hovered) {
                return Theme.surfaceHover
            }
            return Theme.surface
        }
        
        Behavior on color {
            ColorAnimation {
                duration: Theme.animationDurationFast
                easing.type: Theme.animationEasingType
            }
        }
        
        Behavior on border.color {
            ColorAnimation {
                duration: Theme.animationDurationFast
                easing.type: Theme.animationEasingType
            }
        }
        
        // Checkmark
        Text {
            id: checkmark
            anchors.centerIn: parent
            text: FluentIconGlyph.checkMarkGlyph
            font.family: "Segoe Fluent Icons"
            font.pixelSize: 12
            color: Theme.textOnPrimary
            opacity: control.checked ? 1 : 0
            scale: control.checked ? 1 : 0
            
            Behavior on opacity {
                NumberAnimation {
                    duration: Theme.animationDurationFast
                    easing.type: Easing.OutBack
                }
            }
            
            Behavior on scale {
                NumberAnimation {
                    duration: Theme.animationDurationFast
                    easing.type: Easing.OutBack
                }
            }
        }
        
        // Indeterminate state
        Rectangle {
            visible: control.checkState === Qt.PartiallyChecked
            anchors.centerIn: parent
            width: parent.width - 8
            height: 2
            color: Theme.textOnPrimary
            radius: 1
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
