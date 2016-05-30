using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace Guardian
{
    public class GuardianRequest
    {
        public DateTime TimeStamp { get; set; }
        public int Port { get; set; }
        public string Host { get; set; }
    }
}
