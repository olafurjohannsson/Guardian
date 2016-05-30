using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace Guardian
{

    public class GuardianRequestGroupEntity
    {
        public int Count { get; set; }
        public DateTime Started { get; set; }
        public DateTime Ended { get; set; }
        public string Host { get; set; }
    }
}
