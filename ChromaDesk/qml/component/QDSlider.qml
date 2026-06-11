// Fluent Design Slider Component
import QtQuick
import QtQuick.Controls as Controls

Controls.Slider {
    id: control
    
    // ============ Custom Properties ============
    
    property bool showValue: false
    property int decimals: 0
    
    // ============ Size & Style ============
    
    implicitWidth: 200
    implicitHeight: Theme.buttonHeightMedium
    
    // ============ Background (Track) ============
    
    background: Rectangle {
        x: control.leftPadding
        y: control.topPadding + control.availableHeight / 2 - height / 2
        width: control.availableWidth
        height: 4
        radius: 2
        color: Theme.surface
        border.width: Theme.borderWidthThin
        border.color: Theme.border
        
        // Filled portion
        Rectangle {
            width: control.visualPosition * parent.width
            height: parent.height
            radius: parent.radius
            color: control.enabled ? Theme.primary : Theme.primaryDisabled
            
            Behavior on width {
                NumberAnimation {
                    duration: 50
                    easing.type: Easing.OutCubic
                }
            }
        }
    }
    
    // ============ Handle ============
    
    handle: Rectangle {
        x: control.leftPadding + control.visualPosition * (control.availableWidth - width)
        y: control.topPadding + control.availableHeight / 2 - height / 2
        width: 20
        height: 20
        radius: 10
        color: {
            if (!control.enabled) {
                return Theme.primaryDisabled
            }
            if (control.pressed) {
                return Theme.primaryPressed
            }
            if (control.hovered) {
                return Theme.primaryHover
            }
            return Theme.primary
        }
        border.width: 2
        border.color: Theme.textOnPrimary
        
        Behavior on color {
            ColorAnimation {
                duration: Theme.animationDurationFast
                easing.type: Theme.animationEasingType
            }
        }
        
        // Scale effect when pressed
        scale: control.pressed ? 1.1 : 1.0
        
        Behavior on scale {
            NumberAnimation {
                duration: Theme.animationDurationFast
                easing.type: Easing.OutCubic
            }
        }
        
        // Value tooltip
        Rectangle {
            visible: control.showValue && control.pressed
            anchors.bottom: parent.top
            anchors.bottomMargin: 8
            anchors.horizontalCenter: parent.horizontalCenter
            width: valueText.width + Theme.spacingMedium * 2
            height: 24
            radius: Theme.radiusSmall
            color: Theme.surface
            border.width: Theme.borderWidthThin
            border.color: Theme.border
            
            Text {
                id: valueText
                anchors.centerIn: parent
                text: control.value.toFixed(control.decimals)
                font.family: Theme.fontFamily
                font.pixelSize: Theme.fontSizeSmall
                color: Theme.text
            }
        }
    }
    
    // Hover cursor
    HoverHandler {
        cursorShape: control.enabled ? Qt.PointingHandCursor : Qt.ForbiddenCursor
    }
}
