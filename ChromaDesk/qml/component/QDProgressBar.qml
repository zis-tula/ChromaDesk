// Fluent Design ProgressBar Component
import QtQuick
import QtQuick.Controls as Controls

Controls.ProgressBar {
    id: control
    
    // ============ Custom Properties ============
    
    property bool showValue: false
    property color progressColor: Theme.primary
    property bool animated: true
    
    // ============ Size & Style ============
    
    implicitWidth: 200
    implicitHeight: 6
    
    // ============ Background ============
    
    background: Rectangle {
        implicitWidth: control.width
        implicitHeight: control.height
        radius: height / 2
        color: Theme.surface
        border.width: Theme.borderWidthThin
        border.color: Theme.border
    }
    
    // ============ Content (Progress indicator) ============
    
    contentItem: Item {
        implicitWidth: control.width
        implicitHeight: control.height
        
        // Determinate progress
        Rectangle {
            visible: !control.indeterminate
            width: control.visualPosition * parent.width
            height: parent.height
            radius: height / 2
            color: control.progressColor
            
            Behavior on width {
                enabled: control.animated
                NumberAnimation {
                    duration: Theme.animationDurationMedium
                    easing.type: Easing.OutCubic
                }
            }
        }
        
        // Indeterminate progress (animated)
        Rectangle {
            id: indeterminateBar
            visible: control.indeterminate
            width: parent.width * 0.3
            height: parent.height
            radius: height / 2
            color: control.progressColor
            
            SequentialAnimation on x {
                running: control.indeterminate
                loops: Animation.Infinite
                
                NumberAnimation {
                    from: -indeterminateBar.width
                    to: control.width
                    duration: 1500
                    easing.type: Easing.InOutQuad
                }
            }
        }
        
        // Value text
        Text {
            visible: control.showValue && !control.indeterminate
            anchors.centerIn: parent
            text: Math.round(control.value * 100) + "%"
            font.family: Theme.fontFamily
            font.pixelSize: Theme.fontSizeSmall
            color: Theme.text
        }
    }
    
    // Clip the progress bar
    clip: true
}
