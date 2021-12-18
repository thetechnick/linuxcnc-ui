import 'dart:async';

import 'package:angular/angular.dart';

@Component(
  selector: 'lcnc-clock',
  styleUrls: ['clock_component.css'],
  templateUrl: 'clock_component.html',
  directives: [],
  pipes: [DatePipe],
)
class ClockComponent implements OnInit, OnDestroy {
  DateTime now;
  Timer _timer;

  @override
  void ngOnInit() {
    _timer = Timer.periodic(Duration(seconds: 1), (t) {
      now = DateTime.now();
    });
  }

  @override
  void ngOnDestroy() {
    _timer.cancel();
    _timer = null;
  }
}
