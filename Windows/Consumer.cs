using System;
using System.Collections.Generic;
using System.Collections.Concurrent;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace Guardian
{
    public class Consumer<T> : IConsumer<T>
    {
        private BlockingCollection<T> collection;

        public Consumer(BlockingCollection<T> collection)
        {
            this.collection = collection;
        }


        public T Consume()
        {
            if (this.collection == null)
            {
                throw new NullReferenceException(string.Format("Collection is null in {0}", this.GetType().Name));
            }

            while (true)
            {
                T item;
                if (this.collection.TryTake(out item))
                {
                    return item;
                }
            }
        }
    }
}
